// GNU GPL v3 License
// Copyright (c) 2017 github.com:go-trellis

package config

import (
	"math/big"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-trellis/common/formats"
)

const (
	includeReg = `\$\{([0-9|a-z|A-Z]|\.)+\}`
)

// AdapterConfig default config adapter
type AdapterConfig struct {
	ConfigFile   string
	ConfigString string

	readerType ReaderType
	reader     Reader
	locker     sync.RWMutex
	configs    map[string]interface{}
}

// NewAdapterConfig return default config adapter
// name is file's path
func NewAdapterConfig(filepath string) (Config, error) {
	if len(filepath) == 0 {
		return nil, ErrInvalidFilePath
	}
	a := &AdapterConfig{
		ConfigFile: filepath,
		configs:    make(map[string]interface{}),
	}

	err := a.init(OptionFile(filepath))
	if err != nil {
		return nil, err
	}

	return a.copy(), nil
}

func (p *AdapterConfig) init(opts ...Option) (err error) {
	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}

	if len(p.ConfigFile) > 0 {
		p.readerType = fileToReaderType(p.ConfigFile)
		switch p.readerType {
		case ReaderTypeJSON:
			p.reader = NewJSONReader()
		case ReaderTypeYAML:
			p.reader = NewYAMLReader()
		default:
			return ErrNotSupportedReaderType
		}

		err = p.reader.Read(p.ConfigFile, &p.configs)
	}
	if err != nil {
		return ErrValueNil
	}

	if len(p.ConfigString) > 0 {
		switch p.readerType {
		case ReaderTypeJSON:
			p.reader = NewJSONReader()
			err = ParseJSONConfig([]byte(p.ConfigString), &p.configs)
		case ReaderTypeYAML:
			p.reader = NewYAMLReader()
			err = ParseYAMLConfig([]byte(p.ConfigString), &p.configs)
		default:
			return ErrNotSupportedReaderType
		}
	}
	if err != nil {
		return ErrValueNil
	}

	err = p.copyDollarSymbol()

	return
}

// GetKeys get map keys
func (p *AdapterConfig) GetKeys() []string {
	p.locker.RLock()
	defer p.locker.RUnlock()

	var keys []string
	for key := range p.configs {
		keys = append(keys, key)
	}
	return keys
}

func (p *AdapterConfig) copy() *AdapterConfig {
	p.locker.RLock()
	defer p.locker.RUnlock()

	values := DeepCopy(p.configs)

	valuesMap := values.(map[string]interface{})
	return &AdapterConfig{
		ConfigFile:   p.ConfigFile,
		ConfigString: p.ConfigString,
		readerType:   p.readerType,
		reader:       p.reader,
		configs:      valuesMap,
	}
}

// GetTimeDuration return time in p.configs by key
func (p *AdapterConfig) GetTimeDuration(key string, defValue ...time.Duration) time.Duration {
	return formats.ParseStringTime(strings.ToLower(p.GetString(key)))
}

// GetByteSize return time in p.configs by key
func (p *AdapterConfig) GetByteSize(key string) *big.Int {
	return formats.ParseStringByteSize(strings.ToLower(p.GetString(key)))
}

// GetInterface return a interface object in p.configs by key
func (p *AdapterConfig) GetInterface(key string, defValue ...interface{}) (res interface{}) {

	var err error
	var v interface{}

	defer func() {
		if err != nil {
			if len(defValue) == 0 {
				return
			}
			res = defValue[0]
		} else {
			res = v
		}
	}()

	if key == "" {
		return ErrInvalidKey
	}

	v, err = p.getKeyValue(key)
	return
}

// GetString return a string object in p.configs by key
func (p *AdapterConfig) GetString(key string, defValue ...string) (res string) {

	var ok bool
	defer func() {
		if ok || len(defValue) == 0 {
			return
		}
		res = defValue[0]
	}()
	v := p.GetInterface(key, defValue)

	res, ok = v.(string)
	return
}

// GetBoolean return a bool object in p.configs by key
func (p *AdapterConfig) GetBoolean(key string, defValue ...bool) (b bool) {

	var ok bool
	defer func() {
		if ok || len(defValue) == 0 {
			return
		}
		b = defValue[0]
	}()
	v := p.GetInterface(key, defValue)

	switch reflect.TypeOf(v).String() {
	case "bool":
		ok, b = true, v.(bool)
	case "string":
		ok, b = true, (v.(string) == "ON" || v.(string) == "on")
	}

	return
}

// GetInt return a int object in p.configs by key
func (p *AdapterConfig) GetInt(key string, defValue ...int) (res int) {

	var err error
	defer func() {
		if err != nil {
			if len(defValue) == 0 {
				return
			}
			res = defValue[0]
		}
	}()

	v, e := formats.ToInt64(p.GetInterface(key, defValue))
	if e != nil {
		err = e
		return
	}
	return int(v)
}

// GetFloat return a float object in p.configs by key
func (p *AdapterConfig) GetFloat(key string, defValue ...float64) (res float64) {

	var err error
	defer func() {
		if err != nil {
			if len(defValue) == 0 {
				return
			}
			res = defValue[0]
		}
	}()

	v, e := formats.ToFloat64(p.GetInterface(key, defValue))
	if e != nil {
		err = e
		return
	}
	return v
}

// GetList return a list of interface{} in p.configs by key
func (p *AdapterConfig) GetList(key string) (res []interface{}) {

	vS := reflect.Indirect(reflect.ValueOf(p.GetInterface(key)))
	if vS.Kind() != reflect.Slice {
		return nil
	}

	var vs []interface{}
	for i := 0; i < vS.Len(); i++ {
		vs = append(vs, vS.Index(i).Interface())
	}
	return vs
}

// GetStringList return a list of strings in p.configs by key
func (p *AdapterConfig) GetStringList(key string) []string {

	var items []string
	for _, v := range p.GetList(key) {
		item, ok := v.(string)
		if !ok {
			return nil
		}

		items = append(items, item)
	}
	return items
}

// GetBooleanList return a list of booleans in p.configs by key
func (p *AdapterConfig) GetBooleanList(key string) []bool {

	var items []bool
	for _, v := range p.GetList(key) {
		item, ok := v.(bool)
		if !ok {
			return nil
		}

		items = append(items, item)
	}
	return items
}

// GetIntList return a list of ints in p.configs by key
func (p *AdapterConfig) GetIntList(key string) []int {

	var items []int
	for _, v := range p.GetList(key) {
		i, e := formats.ToInt(v)
		if e != nil {
			return nil
		}
		items = append(items, i)
	}
	return items
}

// GetFloatList return a list of floats in p.configs by key
func (p *AdapterConfig) GetFloatList(key string) []float64 {

	var items []float64
	for _, v := range p.GetList(key) {
		f, e := formats.ToFloat64(v)
		if e != nil {
			return nil
		}
		items = append(items, f)
	}
	return items
}

// GetMap get map value
func (p *AdapterConfig) GetMap(key string) Options {

	vm, err := p.getKeyValue(key)
	if err != nil {
		return nil
	}

	mapVM, ok := vm.(map[string]interface{})
	if ok {
		return mapVM
	}

	mapVMs, ok := vm.(map[interface{}]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range mapVMs {
		sk, _ := k.(string)
		result[sk] = v
	}
	return result
}

// GetConfig return object config in p.configs by key
func (p *AdapterConfig) GetConfig(key string) Config {

	vm, err := p.getKeyValue(key)
	if err != nil {
		return nil
	}

	c := &AdapterConfig{
		reader:  p.reader,
		configs: map[string]interface{}{key: vm},
	}

	return c
}

// GetValuesConfig get key's values if values can be Config, or panic
func (p *AdapterConfig) GetValuesConfig(key string) Config {
	opt := p.GetMap(key)
	if opt == nil {
		return nil
	}

	return MapGetter().GenMapConfig(p.readerType, opt)
}

func (p *AdapterConfig) getKeyValue(key string) (vm interface{}, err error) {
	if len(key) == 0 {
		return nil, ErrInvalidKey
	}
	p.locker.RLock()
	defer p.locker.RUnlock()

	switch p.readerType {
	case ReaderTypeJSON:
		return getStringKeyValue(p.configs, key)
	case ReaderTypeYAML:
		return getInterfaceKeyValue(p.configs, key)
	default:
		return nil, ErrNotSupportedReaderType
	}
}

// SetKeyValue set key value into p.configs
func (p *AdapterConfig) SetKeyValue(key string, value interface{}) (err error) {
	if len(key) == 0 {
		return ErrInvalidKey
	}
	p.locker.Lock()
	defer p.locker.Unlock()

	switch p.readerType {
	case ReaderTypeJSON:
		return setStringKeyValue(&p.configs, key, value)
	case ReaderTypeYAML:
		return setInterfaceKeyValue(&p.configs, key, value)
	default:
		return ErrNotSupportedReaderType
	}
}

// Dump return p.configs' bytes
func (p *AdapterConfig) Dump() (bs []byte, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	return p.reader.Dump(p.configs)
}

// Copy return a copy
func (p *AdapterConfig) Copy() Config {
	return p.copy()
}

func (p *AdapterConfig) copyDollarSymbol() error {
	p.locker.RLock()
	defer p.locker.RUnlock()

	switch p.readerType {
	case ReaderTypeJSON:
		return copyJSONDollarSymbol(&p.configs, "", &p.configs)
	case ReaderTypeYAML:
		return copyYAMLDollarSymbol(&p.configs)
	}

	return nil
}

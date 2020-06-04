/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"

	"github.com/go-trellis/config"
	"github.com/spf13/cobra"
	"github.com/ymhhh/tax/handlers"
)

var salariesConfig string

// taxCmd represents the tax command
var taxCmd = &cobra.Command{
	Use:     "tax",
	Aliases: []string{"t"},
	Short:   "计算月工资情况",
	Long: `
计算月工资，并按照五险一金扣除，以及部分可抵扣个税的金额，综合计算按年的月收入情况
./tax t

	完整样例
	./tax --config="tax.yaml" t -c="salaries.yaml"
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始计算个税情况")

		taxes, err := handlers.NewTaxesHandler(cfgFile)
		if err != nil {
			panic(err)
		}

		ss := &handlers.Salaries{}

		if err := config.NewSuffixReader().Read(subCfgFile, ss); err != nil {
			panic(err)
		}

		data, err := taxes.Calc(ss)
		if err != nil {
			panic(err)
		}

		data.Print()
	},
}

var subCfgFile string

func init() {
	rootCmd.AddCommand(taxCmd)

	taxCmd.Flags().StringVarP(&subCfgFile, "subc", "c", "salaries.yaml", "月工资配置文件")
}

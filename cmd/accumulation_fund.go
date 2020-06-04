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
	"log"

	"github.com/go-trellis/config"
	"github.com/ymhhh/tax/handlers"

	"github.com/spf13/cobra"
)

// accumulationFundCmd represents the accumulationFund command
var accumulationFundCmd = &cobra.Command{
	Use:     "accumulationFund",
	Aliases: []string{"fund", "f"},
	Short:   "计算公积金",
	Long: `
通过公积金基数和比例计算公积金
./tax f

	样例:
	./tax --config="tax.yaml" f -c="personal.yaml"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始计算公积金")

		i, err := handlers.NewAccumulationFundHandler(cfgFile)
		if err != nil {
			fmt.Println("读取配置出错")
			return
		}
		var afPersonalInfo handlers.PersonalInfo

		if err := config.NewSuffixReader().Read(accumulationFundConfig, &afPersonalInfo); err != nil {
			log.Fatalln("读取配置失败", err)
			return
		}

		result, err := i.Calc(&afPersonalInfo)
		if err != nil {
			fmt.Println(err)
			return
		}

		result.Print()
	},
}

var accumulationFundConfig string

func init() {
	rootCmd.AddCommand(accumulationFundCmd)

	accumulationFundCmd.Flags().StringVarP(&accumulationFundConfig, "subc", "c", "personal.yaml", "个人信息配置文件路径")
}

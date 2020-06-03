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
	"github.com/ymhhh/tax/handlers"

	"github.com/spf13/cobra"
)

// insuranceCmd represents the insurance command
var insuranceCmd = &cobra.Command{
	Use:     "insurance",
	Aliases: []string{"i"},
	Short:   "计算社会保险",
	Long: `
	通过社会保险基数计算个人的社保缴纳金额

	样例:
	./tax --config="tax.yaml" i -p="personal.yaml"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始计算社会保险")
		i, err := handlers.NewInsurances(cfgFile)
		if err != nil {
			fmt.Println("读取配置文件失败")
			return
		}

		var insuranceInfo handlers.PersonalInfo

		if err := config.NewSuffixReader().Read(insuranceConfig, &insuranceInfo); err != nil {
			fmt.Println("读取配置文件失败")
			return
		}

		result, err := i.Calc(&insuranceInfo)
		if err != nil {
			fmt.Println("计算数据失败")
			return
		}

		result.Print()
	},
}

var insuranceConfig string

func init() {
	rootCmd.AddCommand(insuranceCmd)

	insuranceCmd.Flags().StringVarP(&insuranceConfig, "pconf", "p", "personal.yaml", "个人信息配置文件")
}

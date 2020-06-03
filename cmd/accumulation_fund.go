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

	样例:
	./tax --config="tax.yaml" f -s 1000 -r 12
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始计算公积金")

		i, err := handlers.NewAccumulationFund(cfgFile)
		if err != nil {
			fmt.Println("读取配置出错")
			return
		}

		result, err := i.Calc(afPersonalInfo)
		if err != nil {
			fmt.Println(err)
			return
		}

		result.Print()
	},
}

var afPersonalInfo *handlers.PersonalInfo

func init() {
	rootCmd.AddCommand(accumulationFundCmd)

	afPersonalInfo = &handlers.PersonalInfo{}

	accumulationFundCmd.Flags().Float64VarP(&afPersonalInfo.Salary, "salary", "s", 0, "月薪水，默认为0")
	accumulationFundCmd.Flags().Float64VarP(&afPersonalInfo.AccumulationFundRate, "rate", "r", 12, "缴纳比例")
}

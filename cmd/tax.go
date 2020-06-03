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

// taxCmd represents the tax command
var taxCmd = &cobra.Command{
	Use:     "tax",
	Aliases: []string{"t"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tax called")

		taxes, err := handlers.NewTaxes(cfgFile)
		if err != nil {
			panic(err)
		}

		ss := &handlers.Salaries{}

		if err := config.NewSuffixReader().Read(cfgFile, ss); err != nil {
			panic(err)
		}

		//TODO
		_, err = taxes.Calc(ss)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(taxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// taxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// taxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

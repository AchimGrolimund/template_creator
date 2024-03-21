/*
Copyright © 2024 Achim Grolimund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/AchimGrolimund/template_creator/internal/template"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		_namespace, _ := cmd.Flags().GetString("namespace")
		_type, _ := cmd.Flags().GetString("type")
		_service, _ := cmd.Flags().GetString("service")

		fmt.Printf("Namespace: %s\n", _namespace)
		fmt.Printf("Type: %s\n", _type)
		fmt.Printf("Service: %s\n", _service)

		// Call the CreateTemplate function
		err := template.CreateGitOpsTemplate(_namespace, _type, _service)
		if err != nil {
			fmt.Printf("Failed to create template: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().StringP("namespace", "n", "", "Namespace for the init command")
	initCmd.Flags().StringP("type", "t", "", "Type for the init command")
	initCmd.Flags().StringP("service", "s", "", "Service for the init command")

	initCmd.MarkFlagRequired("namespace")
	initCmd.MarkFlagRequired("type")
	initCmd.MarkFlagRequired("service")
}

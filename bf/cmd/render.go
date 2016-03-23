// Copyright © 2016 Joseph Wright <rjosephwright@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/cloudboss/bossformation"
	"github.com/spf13/cobra"
)

var (
	stackCfg  string
	renderCmd = &cobra.Command{
		Use:   "render",
		Short: "Render a CloudFormation template from a given configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			if stack, err := bf.LoadStackFromFile(stackCfg); err != nil {
				println(err.Error())
			} else {
				if t, err := stack.Render(bf.Context{}); err != nil {
					println(err.Error())
				} else {
					fmt.Printf(t)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&stackCfg, "stack-config", "s", "", "Stack configuration")
}

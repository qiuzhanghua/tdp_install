/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

type SDK struct {
	title   string //SDK 名称
	version string //默认版本
	chip    string //芯片
	os      string //系统
	bit     string //位数

	allowMultiple bool   //允许多个
	port          int    //端口号
	path          string //环境变量设置
	current       bool   //默认版本
	group         string //分组

}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		pwd, err2 := os.Getwd()
		treedir(pwd + "/archives/")
		if err2 != nil {
			fmt.Println("string转map失败", err2)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func treedir(fpath string) {

	// 获取fileinfo
	if finfo, err := os.Stat(fpath); err == nil {
		// 判断是不是目录 如果不是目录而是文件  打印文件path并跳出递归
		if !finfo.IsDir() {
			var filenameWithSuffix string
			filenameWithSuffix = path.Base(fpath) //获取文件名带后缀
			//fmt.Println("filenameWithSuffix =", filenameWithSuffix)
			var fileSuffix string
			fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀
			//fmt.Println("fileSuffix =", fileSuffix)
			var filenameOnly string
			filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名
			fmt.Println(filenameOnly)
			return
		} else {
			// 是目录的情况 打印目录path
			fmt.Println("<========================" + fpath + "========================>")

			f, _ := os.Open(fpath) // 通过目录path open一个file
			defer f.Close()
			names, _ := f.Readdirnames(0) // 通过file的Readdirnames 拿到当前目录下的所有filename
			for _, name := range names {
				newpath := path.Join(fpath, name) // 遍历names 拼接新的fpath
				treedir(newpath)                  // 递归
			}
		}
	}
}

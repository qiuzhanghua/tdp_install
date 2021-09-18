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
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	installCmd.AddCommand(setJdkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.PersistentFlags().StringVarP(&versionStr, "version", "v", "", "version")
}

//解压
func DeCompress(zipFile string, destination string) {

	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(destFile, fileInArchive); err != nil {
			panic(err)
		}

		destFile.Close()
		fileInArchive.Close()
	}
}

var versionStr string //全局变量

/**
 * 获取json字符串中指定字段内容  ioutil.ReadFile()读取字节切片
 * @param    bytes    json字符串字节数组
 * @param    field    可变参数，指定字段
 */
func getJsonField(bytes []byte, field ...string) []byte {
	if len(field) < 1 {
		fmt.Printf("At least two parameters are required.")
		return nil
	}

	//将字节切片映射到指定map上  key：string类型，value：interface{}  类型能存任何数据类型
	var mapObj map[string]interface{}
	json.Unmarshal(bytes, &mapObj)
	var tmpObj interface{}
	tmpObj = mapObj
	for i := 0; i < len(field); i++ {
		tmpObj = tmpObj.(map[string]interface{})[field[i]]
		if tmpObj == nil {
			fmt.Printf("No field specified: %s ", field[i])
			return nil
		}
	}

	result, err := json.Marshal(tmpObj)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return result
}

//一级命令 install
var installCmd = &cobra.Command{
	DisableSuggestions: true,
	Use:                "install",
	Short:              "<candidate> [version] ",
	Example:            "install jdk",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No candidate provided.")
		for _, arg := range args {
			fmt.Println("Arg:", arg)
		}
	},
}

var current_os string //当前操作系统

// 二级子命令 jdk
var setJdkCmd = &cobra.Command{
	DisableSuggestions: true,
	Use:                "jdk",
	Long:               " the SDK to install: jdk eg: $ sdk list",
	Run: func(cmd *cobra.Command, args []string) {
		current_os = runtime.GOOS //获得当前操作系统
		//判断指定安装版本
		if len(versionStr) > 0 {
			fmt.Printf("指定版本[%s]\n", versionStr)
		} else {
			//没有指定版本，选择默认版本安装
			bytes, _ := ioutil.ReadFile("tmp/windows_current.json")
			//获取JDK 安装信息
			str := getJsonField(bytes, "current", "jdk")
			install(str)

		}

	},
}
var pwd string //当前操作系统
//安装软件
func install(str []byte) {

	pwd, err2 := os.Getwd()
	fmt.Println(pwd)
	if err2 != nil {
		fmt.Println("---", err2)
	}
	var resMap map[string]string
	err := json.Unmarshal([]byte(str), &resMap)
	if err != nil {
		fmt.Println("string转map失败", err)
	}
	source := pwd + "/archives/jdk/" + resMap["path"] + ".zip"
	target := pwd + "/candidates/"
	//TODO 暂时做法为解压，日后要替换成软链接

	_, isExistError := os.Lstat(target + resMap["path"])
	fmt.Print(target + resMap["path"])
	if os.IsNotExist(isExistError) {

		DeCompress(source, target)
		err = os.Link(target+resMap["path"], target+"/current/jdk")
		if err != nil {
			fmt.Print("鏈接失效？？")
		}

	} else {
		fmt.Println("软件已安装可选择版本")
	}

	//设置环境变量
	os.Setenv("JAVA_HOME", resMap["env"])
	//fmt.Printf("jdk 环境变量=[%s]\n", os.Getenv("JAVA_HOME"))
	fmt.Print("设置环境变量：" + os.Getenv("JAVA_HOME"))

}

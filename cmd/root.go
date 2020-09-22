// Package cmd provides ...
package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"tosilk/silk"

	"github.com/spf13/cobra"
)

var (
	filePath   string
	fileBase64 string
	fileURL    string
	output     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tosilk",
	Short: "音频文件转成silk格式",
	Long: `
    >>>请保证 ffmpeg 命令可用<<<

    本地使用: 

    1. 本地文件转换(输出base64)
    ./tosilk -f test.mp3

    2. 由base64编码转换(输出base64)
    ./tosilk -b base64-string

    3. 由音频下载链接转换(输出base64)
    ./tosilk -u download-url

    *如果需要直接生成本地文件, 在命令后面指定 -o 参数, 如:
    ./tosilk -f test.mp3 -o test.silk
    (其他的一样)

    开启HTTP服务:
    ./tosilk help server查看帮助
    `,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			b64 string
			err error
		)
		if filePath != "" {
			b64, err = silk.FileToSilkBase64(filePath)
		} else if fileBase64 != "" {
			b64, err = silk.Base64ToSilkBase64(filePath)
		} else if fileURL != "" {
			b64, err = silk.URLToSilkBase64(filePath)
		} else {
			fmt.Println("请使用 tosilk -h 查看帮助")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if output == "" {
			fmt.Print(b64)
		} else {
			data, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				ioutil.WriteFile(output, data, os.ModePerm)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "待转换音频文件路径")
	rootCmd.Flags().StringVarP(&fileBase64, "base64", "b", "", "带转换音频文件的base64编码字符串")
	rootCmd.Flags().StringVarP(&fileURL, "url", "u", "", "待转换音频文件的下载链接")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "输出到文件的路径, 默认输出base64字符串")
}

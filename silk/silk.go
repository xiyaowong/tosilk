// Package silk provides ...
package silk

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
	"tosilk/util"
)

var encoder string

func downloadCodec(url string, path string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, body, os.ModePerm)
	return
}

func init() {
	// fmt.Println("====请确保已下载 FFmpeg 并已设置环境变量====")
	// 检查依赖
	appPath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	codecDir := path.Join(path.Dir(appPath), "codec")
	if !util.FileExist(codecDir) {
		os.MkdirAll(codecDir, os.ModePerm)
	}

	goos := runtime.GOOS
	arch := runtime.GOARCH
	if goos == "windows" && arch == "amd64" {
		encoder = path.Join(codecDir, "windows-encoder.exe")
		if !util.FileExist(encoder) {
			fmt.Println("下载Windows依赖")
			if err = downloadCodec("https://cdn.jsdelivr.net/gh/xiyaowong/tosilk/codec/windows-encoder-exe", encoder); err != nil {
				fmt.Printf("下载依赖失败, %v\n", err)
				os.Exit(1)
			}
		}
	} else if goos == "linux" && arch == "amd64" {
		encoder = path.Join(codecDir, "linux-amd64-encoder")
		if !util.FileExist(encoder) {
			fmt.Println("下载linux amd64依赖")
			if err = downloadCodec("https://cdn.jsdelivr.net/gh/xiyaowong/tosilk/codec/linux-amd64-encoder", encoder); err != nil {
				fmt.Printf("下载依赖失败, %v\n", err)
				os.Exit(1)
			}
		}
	} else if goos == "linux" && arch == "arm64" {
		encoder = path.Join(codecDir, "linux-arm64-encoder")
		if !util.FileExist(encoder) {
			fmt.Println("下载linux arm64依赖")
			if err = downloadCodec("https://cdn.jsdelivr.net/gh/xiyaowong/tosilk/codec/linux-arm64-encoder", encoder); err != nil {
				fmt.Printf("下载依赖失败, %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("%s %s is not supported.\n", goos, arch)
		os.Exit(1)
	}
}

// FileToSilkBase64 音频文件转成silk base64
func FileToSilkBase64(filePath string) (b64 string, err error) {
	// p := fmt.Println
	// 1. 转pcm
	pcmPath := path.Join(path.Dir(filePath), "file.wav")
	defer os.Remove(pcmPath)

	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "s16le", "-ar", "24000", "-ac", "1", "-acodec", "pcm_s16le", pcmPath)
	cmd.Run()
	// 2. 转silk
	silkPath := path.Join(path.Dir(pcmPath), "file.silk")
	defer os.Remove(silkPath)

	cmd = exec.Command(encoder, pcmPath, silkPath, "-quiet", "-tencent")
	cmd.Run()
	// 3. 转base64
	content, err := ioutil.ReadFile(silkPath)
	if err != nil {
		return
	}

	b64 = base64.StdEncoding.EncodeToString(content)
	return
}

// Base64ToSilkBase64 base64 转 silk base64
func Base64ToSilkBase64(b64Str string) (b64 string, err error) {
	data, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return
	}

	tempDir, _ := ioutil.TempDir("", "silk")
	filePath := path.Join(tempDir, "file.wav")
	defer os.RemoveAll(tempDir)

	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return
	}

	return FileToSilkBase64(filePath)
}

// URLToSilkBase64 音频下载链接转base64
func URLToSilkBase64(l string) (b64 string, err error) {
	tempDir, _ := ioutil.TempDir("", "silk")
	filePath := path.Join(tempDir, "file.wav")
	defer os.RemoveAll(tempDir)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, _ := http.NewRequest("GET", l, nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filePath, body, os.ModePerm)
	if err != nil {
		return
	}
	return FileToSilkBase64(filePath)
}

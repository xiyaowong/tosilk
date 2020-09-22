// Package util ...
package util

import (
	"os"
)

// FileExist 检查文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

package test_data

import (
	"fmt"
	"os"
)

// nolint
func findFiles() {
	// 打开文件夹
	dir, err := os.Open("/path/to/folder")
	if err != nil {
		fmt.Println(err)
		return
	}
	// nolint
	defer dir.Close()

	dir2, err := os.Open("/path/to/folder")
	if err != nil {
		fmt.Println(err)
		return
	}
	// nolint
	defer dir2.Close()

	// 读取文件夹中的所有文件
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历文件并打印它们的名称
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

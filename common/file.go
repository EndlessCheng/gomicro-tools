package common

import (
	"io/ioutil"
	"os"
	"path"
	"fmt"
)

// Exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFile(filePath string) *os.File {
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	Check(err)

	if exists(filePath) && !Force {
		filePath += ".gen"
	}

	f, err := os.Create(filePath)
	Check(err)

	return f
}

func ReadText(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("找不到文件:", filePath)
		os.Exit(1)
	}

	return string(data)
}

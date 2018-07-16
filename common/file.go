package common

import (
	"io/ioutil"
	"os"
	"path"
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

	if exists(filePath) {
		filePath += ".gen"
	}

	f, err := os.Create(filePath)
	Check(err)

	return f
}

func ReadText(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	Check(err)

	return string(data)
}

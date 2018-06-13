package common

import (
	"os"
	"path"
	"io/ioutil"
)

func CreateFile(filePath string) *os.File {
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	Check(err)

	f, err := os.Create(filePath)
	Check(err)

	return f
}

func ReadText(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	Check(err)

	return string(data)
}

package profile

import (
	"os"
)

func ReadFile(fileName string) []byte {
	file, err := os.Open(fileName)
	CheckErr(err)
	defer file.Close()

	fileInfo, err := file.Stat()
	CheckErr(err)

	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	CheckErr(err)

	return buf
}

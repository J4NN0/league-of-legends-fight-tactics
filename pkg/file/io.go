package file

import (
	"os"
)

func Create(fileName string) {
	fi, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
}

func Write(fileName, content string) {
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = fo.WriteString(content)
	if err != nil {
		panic(err)
	}
}

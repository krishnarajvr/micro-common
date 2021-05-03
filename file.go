package common

import (
	"bytes"
	"io/ioutil"
	"log"
)

func FileReplaceString(file string, find string, replace string) bool {
	input, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err)
		return false
	}

	output := bytes.Replace(input, []byte(find), []byte(replace), -1)

	if err = ioutil.WriteFile(file, output, 0666); err != nil {
		log.Println(err)
		return false
	}

	return true
}

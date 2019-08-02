package ast_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _, update = os.LookupEnv("update")

func runTestOnFile(output, goldenFileName string) error {
	goldenFileLocation := filepath.Join("testdata", goldenFileName)
	outputBytes := []byte(output)
	if update {
		if err := ioutil.WriteFile(goldenFileLocation, outputBytes, 0644); err != nil {
			return errors.New("failed to update golden file: " + err.Error())
		}
	}
	correct, err := ioutil.ReadFile(goldenFileLocation)
	if err != nil {
		return errors.New("failed to read golden file: " + err.Error())
	}
	if !bytes.Equal(correct, outputBytes) {
		return errors.New("bytes do not match!")
	}
	return nil
}

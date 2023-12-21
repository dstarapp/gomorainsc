package utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(filename string) ([]byte, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	return ioutil.ReadAll(fs)
}

func ReadFileByLine(filename string) ([]string, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	rd := bufio.NewReader(fs)
	ret := make([]string, 0)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			line = strings.TrimSpace(line)
			if line != "" {
				ret = append(ret, line)
			}
			break
		}
		ret = append(ret, strings.TrimSpace(line))
	}

	return ret, nil
}

func IteratorFileByLine(filename string, fn func(line string, index int)) error {
	fs, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()

	bufio.NewScanner(fs)
	rd := bufio.NewReader(fs)
	index := 0
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			line = strings.TrimSpace(line)
			if line != "" {
				fn(strings.TrimSpace(line), index)
			}
			break
		}
		fn(strings.TrimSpace(line), index)
		index++
	}

	return nil
}

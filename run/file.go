package run

import (
	"os"
	"io"
	"path/filepath"
	"fmt"
	"strings"
)

func FileOrDirIsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func walkFileList(path string, action func(file string) (error)) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {
			return err
		}
		if f.IsDir() {
			return nil
		}
		return action(f.Name())
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func CheckMatch(fileName string, pattern string) bool {

	fileName = strings.TrimSpace(fileName)
	pattern = strings.TrimSpace(pattern)
	if pattern == "*" || pattern == fileName {
		return true;
	}
	prefixMatch := strings.HasPrefix(pattern, "*") && strings.HasSuffix(fileName, pattern[1:])
	suffixMatch := strings.HasSuffix(pattern, "*") && strings.HasPrefix(fileName, pattern[0:len(pattern)-1])
	if prefixMatch || suffixMatch {
		return true
	}
	return false;
}

func CheckMatchInFilePattern(fileName string, patterns []string) bool {
	for _, p := range patterns {
		if CheckMatch(fileName, p) {
			return true
		}
	}
	return false
}

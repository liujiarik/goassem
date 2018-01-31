package run

import (
	"os"
	"io"
	"github.com/pierrre/archivefile/zip"
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

func gzipCompress(inFilePath string, outFilePath string) error {
	return zip.ArchiveFile(inFilePath, outFilePath, nil)
}

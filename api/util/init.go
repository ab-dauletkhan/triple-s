package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

func createFileWithDefaultContent(filePath string, header []string) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, core.FilePerm)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	fileStat, err := f.Stat()
	if err != nil {
		log.Println(err)
		return err
	}

	if fileStat.Size() == 0 {
		_, err := f.WriteString(strings.Join(header, ",") + "\n")
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func InitDir() error {
	err := os.MkdirAll(core.Dir, core.DirPerm)
	if err != nil {
		return err
	}

	return createFileWithDefaultContent(filepath.Join(core.Dir, core.BucketsFile), core.BucketsCSVHeader)
}

func InitObjectFile(bucketName string) error {
	err := os.MkdirAll(filepath.Join(core.Dir, bucketName), core.DirPerm)
	if err != nil {
		return err
	}

	return createFileWithDefaultContent(filepath.Join(core.Dir, bucketName, core.ObjectsFile), core.ObjectsCSVHeader)
}

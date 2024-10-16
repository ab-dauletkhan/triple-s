package util

import (
	"os"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

func InitDir() error {
	err := os.MkdirAll(core.Dir, 0o755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(core.Dir+"/buckets.xml", os.O_CREATE|os.O_RDWR, 0o766)
	if err != nil {
		return err
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		return err
	}

	if fs.Size() == 0 {
		_, err := f.WriteString("<Buckets></Buckets>")
		if err != nil {
			return err
		}
	}

	return nil
}

func InitObjectFile(bucketName string) error {
	err := os.MkdirAll(core.Dir+"/"+bucketName, 0o755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(core.Dir+"/"+bucketName+"/objects.xml", os.O_CREATE|os.O_RDWR, 0o766)
	if err != nil {
		return err
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		return err
	}

	if fs.Size() == 0 {
		_, err := f.WriteString("<Objects></Objects>")
		if err != nil {
			return err
		}
	}

	return nil
}

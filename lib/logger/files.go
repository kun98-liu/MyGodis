package logger

import (
	"fmt"
	"os"
)

func mustOpen(filename string, dir string) (*os.File, error) {
	perm := checkPermission(dir)
	if perm {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	err := isNotExistMkDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return f, nil

}

func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func isNotExistMkDir(src string) error {
	if notExist := checkNotExist(src); notExist {
		if err := mkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}
func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

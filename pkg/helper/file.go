package helper

import (
	"github.com/pkg/errors"
	"os"
)

// 파일이나 디렉토리가 존재하는지 확인
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

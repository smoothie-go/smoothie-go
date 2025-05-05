package temp

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
)

var (
	tempFiles     []string = nil
	tempDirectory string   = ""
)

func InitTemp(args *cli.Arguments) error {
	tempDirectory = portable.GetTempPath(filepath.Base(args.InputFile))
	if _, err := os.Stat(tempDirectory); !os.IsNotExist(err) {
		return nil
	}
	err := os.Mkdir(tempDirectory, 0755)
	if err != nil {
		return err
	}
	return nil
}

func Join(fileName string) (string, error) {
	if tempDirectory == "" {
		return "", errors.New("Dev error: Initiate tempDirectory string before using CreateTempFile()")
	}
	return filepath.Join(tempDirectory, fileName), nil
}

func RegisterTempFile(fileName string) error {
	if tempDirectory == "" {
		return errors.New("Dev error: Initiate tempDirectory string before using CreateTempFile()")
	}
	if _, err := os.Stat(filepath.Join(tempDirectory, fileName)); os.IsNotExist(err) {
		return errors.New("Error: Unable to register temp file, " + err.Error())
	}
	tempFiles = append(tempFiles, fileName)
	return nil
}

func DeleteTempFiles() error {
	if tempDirectory == "" || tempFiles == nil || len(tempFiles) == 0 {
		return errors.New("No temp directory or no temp files")
	}

	for _, file := range tempFiles {
		err := os.Remove(filepath.Join(tempDirectory, file))
		if err != nil {
			return err
		}
	}
	err := os.Remove(tempDirectory)
	if err == syscall.ENOTEMPTY {
		return errors.New("tempDirectory is not empty, maybe you forgot to register a temp file")
	} else if err != nil {
		return err
	}

	return nil
}

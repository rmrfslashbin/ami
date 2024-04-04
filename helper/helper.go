package helper

import (
	"os"
	"path/filepath"
)

type SaveFileInput struct {
	Filename string
	Data     []byte
	FileMode *os.FileMode
}

func SaveFile(input *SaveFileInput) error {
	if input.FileMode == nil {
		input.FileMode = new(os.FileMode)
		*input.FileMode = 0644
	}
	fqpn, err := filepath.Abs(input.Filename)
	if err != nil {
		return err
	}

	return os.WriteFile(fqpn, input.Data, *input.FileMode)
}

func LoadFile(filename string) ([]byte, error) {
	fqpn, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fqpn)
}

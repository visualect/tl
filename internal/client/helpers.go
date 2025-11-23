package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/visualect/tl/internal/dto"
)

func SaveFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0600)
}

func IsFileExists(filename string) ([]byte, bool) {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false
		}
	}
	return file, true
}

func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func GetToken(filename string) (string, error) {
	file, ok := IsFileExists(filename)
	if !ok {
		return "", fmt.Errorf("error: get token from %s", filename)
	}
	var t dto.LoginResponse
	err := json.Unmarshal(file, &t)
	if err != nil {
		return "", fmt.Errorf("error: get token from %s", filename)
	}
	return t.Token, nil
}

func GetMessage(b []byte) (string, error) {
	type M struct {
		Message string `json:"message"`
	}
	var m M
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&m)
	if err != nil {
		return "", err
	}
	return m.Message, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/visualect/tl/internal/dto"
)

var backendURL = "http://localhost:8000"

func Login(login string, password string) ([]byte, error) {
	creds := dto.LoginUserRequest{Login: login, Password: password}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(creds)
	if err != nil {
		return nil, err
	}

	data, code, err := baseRequest(backendURL+"/login", "POST", &b)
	if code != http.StatusOK {
		return nil, errors.New("failed to log in")
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func SignUp(login string, password string) ([]byte, error) {
	creds := dto.RegisterUserRequest{Login: login, Password: password}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(creds)
	if err != nil {
		return nil, err
	}

	data, code, err := baseRequest(backendURL+"/signup", "POST", &b)
	if code != http.StatusOK {
		return nil, errors.New("failed to log in")
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func AddTask(task string) error {
	newTask := dto.AddTaskRequest{Task: task}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(newTask)
	if err != nil {
		return err
	}

	data, code, err := baseRequest(backendURL+"/tasks", "POST", &b)
	if code < 200 || code >= 300 {
		return errors.New("failed to add task")
	}

	return nil

}

func baseRequest(url, method string, body io.Reader) (data []byte, statusCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, resp.StatusCode, nil
}

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

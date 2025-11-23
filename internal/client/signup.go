package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/visualect/tl/internal/dto"
)

func SignUp(login string, password string) (string, error) {
	creds := dto.RegisterUserRequest{Login: login, Password: password}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(creds)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", os.Getenv("BACKEND_URL")+"/signup", &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		msg, _ := GetMessage(data)
		return "", errors.New(msg)
	}

	var l dto.RegisterResponse
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&l)

	return l.Login, nil
}

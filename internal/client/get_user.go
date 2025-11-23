package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/visualect/tl/internal/dto"
)

func GetUser() (dto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", os.Getenv("BACKEND_URL")+"/me", nil)
	if err != nil {
		return dto.UserResponse{}, err
	}

	token, err := GetToken(os.Getenv("AUTH_FILENAME"))
	if err != nil {
		return dto.UserResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.UserResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		msg, _ := GetMessage(data)
		return dto.UserResponse{}, errors.New(msg)
	}

	var u dto.UserResponse
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&u)
	if err != nil {
		log.Fatal(err)
	}

	return u, nil
}

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

	"github.com/visualect/tl/internal/models"
)

func GetTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", os.Getenv("BACKEND_URL")+"/tasks", nil)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(os.Getenv("AUTH_FILENAME"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		msg, _ := GetMessage(data)
		return nil, errors.New(msg)
	}

	var list []models.Task
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

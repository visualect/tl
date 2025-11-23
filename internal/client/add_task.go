package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/visualect/tl/internal/dto"
)

func AddTask(task string) ([]byte, error) {
	newTask := dto.AddTaskRequest{Task: task}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(newTask)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", backendURL+"/tasks", &b)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(AuthFilename)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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

	return data, nil
}

package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

func DeleteTask(id string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "DELETE", backendURL+"/tasks/"+id, nil)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(AuthFilename)
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

	return data, nil
}

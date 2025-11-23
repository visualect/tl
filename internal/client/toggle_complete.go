package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ToggleCompleteTask(localID int) error {
	tasks, err := GetTasks()
	if err != nil {
		return err
	}
	var serverID int
	for i, t := range tasks {
		if i+1 == localID {
			serverID = t.ID
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "PATCH", os.Getenv("BACKEND_URL")+"/tasks/"+strconv.Itoa(serverID), nil)
	if err != nil {
		return err
	}

	token, err := GetToken(os.Getenv("AUTH_FILENAME"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		msg, _ := GetMessage(data)
		return errors.New(msg)
	}

	return nil
}

package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/piger/hermano/internal/config"
)

const endpoint = "https://api.pushover.net/1/messages.json"

type Payload struct {
	Token   string `json:"token"`
	UserKey string `json:"user"`
	Message string `json:"message"`
}

func Notify(config *config.Config, message string) error {
	p := Payload{
		Token:   config.APIToken,
		UserKey: config.UserKey,
		Message: message,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		return err
	}

	resp, err := http.Post(endpoint, "application/json", &buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error sending notification: status=%d", resp.StatusCode)
	}

	return nil
}

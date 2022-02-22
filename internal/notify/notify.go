package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/piger/hermano/internal/config"
)

// Pushover's endpoint to send Push Notifications.
const endpoint = "https://api.pushover.net/1/messages.json"

// payload is the data structure used to send a Push notification in Pushover.
type payload struct {
	Token   string `json:"token"`
	UserKey string `json:"user"`
	Message string `json:"message"`
}

// Notify sends a Push notification with Pushvoer.
func Notify(config *config.Config, message string) error {
	p := payload{
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error sending notification: status=%d", resp.StatusCode)
	}

	return nil
}

package pushover

import (
	"fmt"

	"github.com/gregdel/pushover"
	psh "github.com/gregdel/pushover"
)

func Notify(app *psh.Pushover, recipient *psh.Recipient, title, body string) error {
	message := pushover.NewMessageWithTitle(body, title)

	_, err := app.SendMessage(message, recipient)
	if err != nil {
		return fmt.Errorf("failed to send pushover notification: %s", err)
	}
	return nil
}

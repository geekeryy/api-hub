package email_test

import (
	"os"
	"testing"

	"github.com/geekeryy/api-hub/core/email"
)

func TestEmail(t *testing.T) {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	if len(apiKey) == 0 {
		t.Skip("MAILGUN_API_KEY is not set")
	}
	err := email.New("mailgun.jiangyang.online", apiKey).SendMailGun(&email.SendMsg{
		Subject: "test",
		Body:    "test",
		To:      []string{"jiangyang.me@gmail.com"},
	}, "api-hub@mailgun.jiangyang.online")
	if err != nil {
		t.Fatal(err)
	}
}

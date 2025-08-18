package facebook_test

import (
	"context"
	"os"
	"testing"

	"github.com/geekeryy/api-hub/core/facebook"
)

func TestGetUserInfo(t *testing.T) {
	appID := os.Getenv("FACEBOOK_APP_ID")
	if len(appID) == 0 {
		t.Skip("FACEBOOK_APP_ID is not set")
	}
	apiKey := os.Getenv("FACEBOOK_APP_SECRET")
	if len(apiKey) == 0 {
		t.Skip("FACEBOOK_APP_SECRET is not set")
	}
	app := facebook.NewFaceBookApp(appID, apiKey)
	user, err := app.GetUserInfo(context.Background(), "EAANFsT4aXNUBO4jBNY0ZAFveOap3sqJHH5MlCSwQnZC1X5Qpi8oGvkol0iyULGSK8U1O8ZB57NnnCygeILfazltQTZC2uNvOZCFz7wK0yF5FDbbLyJFCp2tutdMD5kZCH7mziZAlDtCM6MID8yZBlMg8IfZCLsti9IGyWPf8wwPuxVyjXwKUSd5ShyqGs6oKfxT5JMLEu0FqfKWPJaH6CxwZDZD")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

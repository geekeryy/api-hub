package facebook_test

import (
	"context"
	"testing"

	"github.com/geekeryy/api-hub/core/facebook"
)

func TestGetUserInfo(t *testing.T) {
	app := facebook.NewFaceBookApp("1345054219852807", "6e134b324739e274bcc07bea608dda2b")
	user, err := app.GetUserInfo(context.Background(), "EAANFsT4aXNUBO4jBNY0ZAFveOap3sqJHH5MlCSwQnZC1X5Qpi8oGvkol0iyULGSK8U1O8ZB57NnnCygeILfazltQTZC2uNvOZCFz7wK0yF5FDbbLyJFCp2tutdMD5kZCH7mziZAlDtCM6MID8yZBlMg8IfZCLsti9IGyWPf8wwPuxVyjXwKUSd5ShyqGs6oKfxT5JMLEu0FqfKWPJaH6CxwZDZD")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

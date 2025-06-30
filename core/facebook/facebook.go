package facebook

import (
	"context"
	"encoding/json"
	"fmt"

	fb "github.com/huandu/facebook/v2"
)

// https://developers.facebook.com/docs/facebook-login/web/js-example
// https://developers.facebook.com/docs/facebook-login/
// https://developers.facebook.com/docs/facebook-login/guides/advanced/manual-flow#checktoken

type UserInfo struct {
	AppID       string   `json:"app_id"`
	Type        string   `json:"type"`
	Application string   `json:"application"`
	ExpiresAt   string   `json:"expires_at"`
	IsValid     bool     `json:"is_valid"`
	IssuedAt    string   `json:"issued_at"`
	UserID      string   `json:"user_id"`
	Scopes      []string `json:"scopes"`
	Email       string   `json:"email"`
}

const (
	FACEBOOK_APP_ID     = "921052483574997"
	FACEBOOK_APP_SECRET = "b7c232f9db402613b09b7c3d13470864"
)

type FaceBookApp struct {
	AppID     string
	AppSecret string
}

func NewFaceBookApp(appID, appSecret string) *FaceBookApp {
	return &FaceBookApp{
		AppID:     appID,
		AppSecret: appSecret,
	}
}

func (f *FaceBookApp) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	app := fb.New(f.AppID, f.AppSecret)
	session := app.Session(accessToken)

	result, err := session.Inspect()
	if err != nil {
		return nil, err
	}
	fmt.Printf("result: %v\n", result)

	var expiresAt, issuedAt string
	var scopes []string
	if result.Get("expires_at") != nil {
		expiresAt = result.Get("expires_at").(json.Number).String()
	}
	if result.Get("issued_at") != nil {
		issuedAt = result.Get("issued_at").(json.Number).String()
	}
	if result.Get("scopes") != nil {
		scopesResult := result.Get("scopes").([]interface{})
		for _, scope := range scopesResult {
			scopes = append(scopes, scope.(string))
		}
	}

	return &UserInfo{
		AppID:       result.Get("app_id").(string),
		Type:        result.Get("type").(string),
		Application: result.Get("application").(string),
		ExpiresAt:   expiresAt,
		IsValid:     result.Get("is_valid").(bool),
		IssuedAt:    issuedAt,
		UserID:      result.Get("user_id").(string),
		Scopes:      scopes,
	}, nil
}

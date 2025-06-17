// Package google @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/9/13 09:49
package google

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

//	{
//	 "sub": "116634423526349493532",
//	 "name": "xingzu quan",
//	 "given_name": "xingzu",
//	 "family_name": "quan",
//	 "picture": "https://lh3.googleusercontent.com/a/ACg8ocKQdbl6qxh_GSppcue1LIy2H5kv8FU5UrwOsSI0Jq3s4N8LkLQ\u003ds96-c",
//	 "email": "quanxingzu@gmail.com",
//	 "email_verified": true
//	}
func GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	res := UserInfo{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil

}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

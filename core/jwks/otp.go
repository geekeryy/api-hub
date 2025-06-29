package jwks

import (
	"bytes"
	"image/png"

	"github.com/pquerna/otp/totp"
)

func GenerateOTP(issuer, account string) (string, []byte, error) {
	// 创建TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
		Period:      30,
	})
	if err != nil {
		return "", nil, err
	}

	// 生成QR码
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", nil, err
	}
	png.Encode(&buf, img)
	return key.Secret(), buf.Bytes(), nil
}

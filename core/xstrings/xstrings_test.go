package xstrings_test

import (
	"testing"

	"github.com/geekeryy/api-hub/core/xstrings"
)

func TestAesCbcEncryptBase64(t *testing.T) {
	plainText := "This is a secret message."
	// secretKey := "1234567890123456"
	cipherText, err := xstrings.AesCbcEncryptBase64(plainText, "public_key_secre", nil)
	if err != nil {
		t.Fatalf("AesCbcEncrypt failed: %v", err)
	}
	t.Logf("cipherText: %s", cipherText)
	decryptText, err := xstrings.AesCbcDecryptBase64(cipherText, "public_key_secre", nil)
	if err != nil {
		t.Fatalf("AesCbcDecrypt failed: %v", err)
	}
	t.Logf("decryptText: %s", decryptText)
}

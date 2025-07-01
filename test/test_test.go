package test_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/core/xcontext"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/gin-gonic/gin"
)

func TestSendEmail(t *testing.T) {
	loador := sync.Map{}
	l, ok := loador.LoadOrStore("test", "v")
	fmt.Println(l, ok)
	fmt.Println(loador.Load("test"))
}

func Test_jwks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create a cryptographic key.
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate given key.\nError: %s", err)
	}

	// Turn the key into a JWK.
	marshalOptions := jwkset.JWKMarshalOptions{
		Private: true,
	}

	metadata := jwkset.JWKMetadataOptions{
		KID: time.Now().Format(time.RFC3339),
	}
	options := jwkset.JWKOptions{
		Marshal:  marshalOptions,
		Metadata: metadata,
	}
	jwk, err := jwkset.NewJWKFromKey(pub, options)
	if err != nil {
		log.Fatalf("Failed to create a JWK from the given key.\nError: %s", err)
	}

	// Write the JWK to the server's storage.
	serverStore := jwkset.NewMemoryStorage()
	err = serverStore.KeyWrite(ctx, jwk)
	if err != nil {
		log.Fatalf("Failed to write the JWK to the server's storage.\nError: %s", err)
	}

	rawJWKS, err := serverStore.JSONPublic(ctx)
	if err != nil {
		log.Fatalf("Failed to get the server's JWKS.\nError: %s", err)
	}
	fmt.Println(string(rawJWKS))

}

func TestValidate(t *testing.T) {
	go server()
	time.Sleep(1 * time.Second)
	params := url.Values{}
	params.Add("id", "05")
	params.Add("id", "25")
	params.Add("id", "35")
	req, err := http.NewRequest("GET", "http://localhost:8080/test?"+params.Encode(), nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("accept-language", "en")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

}

type TestRequest struct {
	// ID   []string `form:"id" validate:"required,dive,oneof=0 1 2 3" comment:"FIELD_USERNAME"`
	Name string `json:"name" comment:"FIELD_USERNAME" validate:"chinese_name" `
}

func server() {
	validator := validate.New([]validate.ValidatorFn{validator.ChineseNameValidator}, []string{"zh", "en"})
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		var req TestRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		if len(c.Request.Header.Get("accept-language")) > 0 {
			ctx = xcontext.WithLang(ctx, c.Request.Header.Get("accept-language"))
		}
		if err := validator.ValidateStruct(ctx, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	g.Run(":8080")
}

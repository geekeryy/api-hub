package test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/geekeryy/api-hub/core/consts"
	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/gin-gonic/gin"
)

func Test_test(t *testing.T) {
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
	// req.Header.Set("accept-language", "")
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
		if len(c.Request.Header.Get("accept-language"))>0{
			ctx=context.WithValue(ctx, consts.AcceptLanguage, c.Request.Header.Get("accept-language"))
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

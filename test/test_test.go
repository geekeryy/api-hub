package test

import (
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/geekeryy/api-hub/core/validate"
	"github.com/geekeryy/api-hub/library/validator"
	"github.com/gin-gonic/gin"
)

func Test_test(t *testing.T) {
	validate.Register([]validate.ValidatorFn{validator.ChineseNameValidator}, []string{"zh","en"})
	go server()
	time.Sleep(1 * time.Second)
	params := url.Values{}
	params.Add("id", "05")
	params.Add("id", "25")
	params.Add("id", "35")
	resp, err := http.Get("http://localhost:8080/test?" + params.Encode())
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
	ID   []string `form:"id" binding:"required,dive,oneof=0 1 2 3" comment:"测试"`
	Name string   `json:"name" comment:"测试" binding:"chinese_name" `
}

func server() {
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		var req TestRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	g.Run(":8080")
}

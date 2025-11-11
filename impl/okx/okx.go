package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/utils"
)

type okx struct {
	apiKey     string
	secretKey  string
	passphrase string
	baseURL    string
	client     *http.Client
}

func NewOKX(apiKey, secretKey string, passphrase string) exchange.Exchange {
	return &okx{
		apiKey:     apiKey,
		secretKey:  secretKey,
		passphrase: passphrase,
		baseURL:    "https://www.okx.com",
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// 生成认证请求
func (o *okx) authRequest(method, requestPath string, body map[string]string) (json.RawMessage, error) {
	method = strings.ToUpper(method)
	var resp okxResp
	var err error
	var bodyString string
	if body == nil {
		bodyString = ""
	}

	client := utils.NewHTTPClient(o.baseURL).SetHeaders(map[string]string{
		"OK-ACCESS-KEY":        o.apiKey,
		"OK-ACCESS-SIGN":       o.generateSignature(method, requestPath, bodyString),
		"OK-ACCESS-TIMESTAMP":  time.Now().UTC().Format(time.RFC3339),
		"OK-ACCESS-PASSPHRASE": o.passphrase,
		"Content-Type":         "application/json",
	})
	switch method {
	case "POST":
	default:
		err = client.Get(requestPath, body).JSON(&resp)
	}
	if err != nil {
		return nil, err
	}
	if resp.Code != "0" {
		return nil, fmt.Errorf("API 返回错误: code=%s, msg=%s", resp.Code, resp.Msg)
	}
	return resp.Data, nil
}

// 生成签名
func (o *okx) generateSignature(method, requestPath string, body string) string {
	// 生成时间戳（ISO 8601 格式）
	timestamp := time.Now().UTC().Format(time.RFC3339)

	message := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(o.secretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

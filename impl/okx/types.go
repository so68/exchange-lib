package okx

import "encoding/json"

type okxResp struct {
	Code string          `json:"code"` // 0 成功
	Msg  string          `json:"msg"`  // 错误信息
	Data json.RawMessage `json:"data"` // 数据
}

package internal

import (
	"encoding/json"
	"log"
)

type HTTPResp struct {
    Message string
}

func (r *HTTPResp) ToBytes() ([]byte) {
    bytes, err := json.Marshal(r)
    if err != nil {
        log.Fatal("数据转换失败")
        return nil
    }
    return bytes
}
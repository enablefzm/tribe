package jsondb

import (
	"encoding/json"
)

func JsdbJson(source interface{}) (string, error) {
	if btVal, err := json.Marshal(source); err != nil {
		return "", err
	} else {
		return string(btVal), nil
	}
}

func JsdbUnJson(strJson string, sTypeOB interface{}) interface{} {
	json.Unmarshal([]byte(strJson), &sTypeOB)
	return sTypeOB
}

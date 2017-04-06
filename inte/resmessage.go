package inte

import (
	"encoding/json"
)

func fInfo(resMap map[string]interface{}) string {
	v, _ := json.Marshal(resMap)
	return string(v)
}

// JSON格式返回消息
func ResMessageCMD(msg string) string {
	mapInfo := GetResMessage("ShowMSG", msg)
	return fInfo(mapInfo)
}

func GetResMessage(cmd string, info interface{}) map[string]interface{} {
	obRes := NewResMessage(cmd)
	obRes.SetInfo(info)
	return obRes.GetRes()
}

func NewResMessageInfo(cmd string) (*ResMessage, map[string]interface{}) {
	obRes := NewResMessage(cmd)
	info := make(map[string]interface{})
	obRes.SetInfo(info)
	return obRes, info
}

func NewResMessage(cmd string) *ResMessage {
	obRes := &ResMessage{
		RES: make(map[string]interface{}),
	}
	obRes.addInfo("CMD", cmd)
	return obRes
}

// 返回信息对象
type ResMessage struct {
	RES map[string]interface{}
}

func (t *ResMessage) SetInfo(val interface{}) {
	t.addInfo("INFO", val)
}

func (t *ResMessage) SetRes(blnRes bool, msg string) {
	t.addInfo("RES", blnRes)
	t.addInfo("MSG", msg)
}

func (t *ResMessage) addInfo(key string, val interface{}) {
	t.RES[key] = val
}

func (t *ResMessage) GetRes() map[string]interface{} {
	return t.RES
}

func (t *ResMessage) GetString() string {
	return fInfo(t.RES)
}

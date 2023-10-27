package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	DingURL    = ""
	DingSercet = ""
)

type DingTalkLogItem struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type DingTalkLog struct {
	Title string                 `json:"title,omitempty"`
	Kvs   map[string]interface{} `json:"items,omitempty"`
}

type DingTalkData struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

// SendDingTalk : 推送重要信息到钉钉
func SendDingTalk(msg *DingTalkLog) error {
	var data = &DingTalkData{MsgType: "markdown"}
	data.Markdown.Title = msg.Title
	data.Markdown.Text = "#### " + msg.Title + "\n"
	for k, v := range msg.Kvs {
		data.Markdown.Text += ">- **" + k + ":** " + fmt.Sprint(v) + " \n"
	}
	return sendDingTalk(DingURL, data)
}

func sendDingTalk(url string, data interface{}) error {
	if url == "" {
		return nil
	}
	t := fmt.Sprint(NowNano() / 1000_000)
	url += fmt.Sprintf("&timestamp=%v&sign=%v", t, SignSHA256(t))
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bts))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{Timeout: time.Second * 5}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	buff, _ := ioutil.ReadAll(rsp.Body)

	log.Println(string(buff))

	return nil
}

func SignSHA256(str string) string {
	str = str + "\n" + DingSercet
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(DingSercet))
	// Write Data to it
	h.Write([]byte(str))
	// Get result and encode as hexadecimal string
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
	// return fmt.Sprintf("%X", h.Sum(nil))
}

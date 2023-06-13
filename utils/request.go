package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tcp-http-golang/global"
)

// HttpPost post url参数方式
func HttpPost(postUrl string) {
	resp, err := http.Post(postUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader("username=test&password=ab123123"))
	if err != nil {
		global.Logger.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	global.Logger.Infof(string(body))
}

// HttpPostForm post Form方式
func HttpPostForm(val string, postUrl string) {
	resp, err := http.PostForm(postUrl,
		url.Values{"aa": {val}})
	//url.Values{"aa": {"{ password: 'auto123123', username: 'auto' }"}})
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	global.Logger.Infof(string(body))
}

// HttpPostJson post json方式
func HttpPostJson(bodyjson string, postUrl string) []byte {

	buffer := bytes.NewBuffer([]byte(bodyjson)) //转换二进制

	//POST请求
	_url := postUrl
	req, err := http.NewRequest("POST", _url, buffer)
	if err != nil {
		global.Logger.Errorf("创建连接出错%", err)
		return []byte("err")
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		global.Logger.Errorf("URL连接出错%s", err)
		return []byte("err")
	}

	defer response.Body.Close()
	//statusCode := response.StatusCode
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Printf("statusCode:%v : %v \n", statusCode, string(body))
	//fmt.Printf("body:%v,%v", string(body))
	return body
}

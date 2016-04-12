package main

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"github.com/Alexniver/logger4go"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

func main() {
	timer := time.NewTicker(24 * time.Hour) //24小时执行一次
	for {
		select {
		case <-timer.C:
			go sign()
		}
	}
}

func sign() {
	logger := logger4go.GetDefaultLogger()

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		logger.Error(err)
		return
	}

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}
	client.Get("http://www.zimuzu.tv/user/login")

	// login form data
	formValue := url.Values{}
	config, err := ParseConfig("config.json")
	if err != nil {
		logger.Error(err)
		return
	}
	for key, val := range config {
		formValue.Set(key, val)
	}

	if err != nil {
		logger.Error(err)
		return
	}
	req, _ := http.NewRequest("POST", "http://www.zimuzu.tv/User/Login/ajaxLogin", bytes.NewBufferString(formValue.Encode()))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	/*content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Error(err)
	}
	fmt.Println("login : ", string(content))*/

	resp, err = client.Get("http://www.zimuzu.tv/user/sign")
	/*content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("get page: ", err)
	}
	fmt.Println(string(content))*/

	time.Sleep(3 * time.Second)
	resp, err = client.Get("http://www.zimuzu.tv/user/login/getCurUserTopInfo")

	/*content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("get Cur User Top Info : ", err)
	}

	fmt.Println(string(content))*/
	resp, err = client.Get("http://www.zimuzu.tv/user/&")

	return //已经不需要再签到了,只要登陆即可

	time.Sleep(15 * time.Second)
	resp, err = client.Get("http://www.zimuzu.tv/user/sign/dosign")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(string(body))
}

func ParseConfig(path string) (result map[string]string, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return
}

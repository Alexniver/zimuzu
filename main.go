package main

import (
	"bytes"
	"encoding/json"
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
	for {
		sign()
		time.Sleep(24 * time.Hour) //24小时执行一次
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

	resp, err = client.Get("http://www.zimuzu.tv/user/sign")

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

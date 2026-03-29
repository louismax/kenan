package kTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// HTTPGet GET请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

// HTTPGetComplex GET请求
func HTTPGetComplex(uri string, param map[string]interface{}, header map[string]string) ([]byte, error) {
	client := &http.Client{}

	if len(param) > 0 {
		uri += `?`
		inx := 1
		for k, v := range param {
			if inx != 1 {
				uri += `&`
			}
			if reflect.TypeOf(v).Kind() == reflect.Map {
				b, _ := json.Marshal(v)
				uri += fmt.Sprintf("%v=%v", k, url.QueryEscape(string(b)))
			} else {
				uri += fmt.Sprintf("%v=%v", k, v)
			}
			inx++
		}
	}

	fmt.Println(uri)
	//提交请求
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	for key, val := range header {
		request.Header.Add(key, val)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	respX, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return respX, nil
}

// HTTPPostJson post请求 JSON
func HTTPPostJson(url string, data interface{}, header map[string]string) ([]byte, error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	//提交请求
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	//request.Header.Add("content-type", "application/json")
	//request.Header.Add("x-aob-signature", signStr)
	//add方式新增Header会被强制转成驼峰
	request.Header.Add("Content-Type", "application/json")
	for key, val := range header {
		request.Header.Add(key, val)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	respX, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return respX, nil
}

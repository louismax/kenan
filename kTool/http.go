package kTool

import (
	"fmt"
	"io"
	"net/http"
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

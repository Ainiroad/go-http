package go_http

import (
	"time"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

func doRequest(req *http.Request, timeout time.Duration, response interface{}) ([]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if response != nil {
		if err := json.Unmarshal(data, response); err != nil {
			return data, err
		}
	}
	return data, nil
}

func Request(addr string, method string, contentType string, body []byte, timeout time.Duration, retry int, response interface{}) ([]byte, error) {
	req, err := http.NewRequest(method, addr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	var (
		rawData []byte
	)
	for i := 0; i < retry; i++ {
		rawData, err = doRequest(req, timeout, response)
		if err == nil {
			break
		}
	}

	return rawData, err
}

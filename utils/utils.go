package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// ReadUnmarshalJSON 简单封装
func ReadUnmarshalJSON(r io.Reader, js interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, js)
}

// HTTPGetJSON 简单封装
func HTTPGetJSON(url string, js interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return ReadUnmarshalJSON(resp.Body, js)
}

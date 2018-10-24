package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// ReadUnmarshalJSON 简单封装
func ReadUnmarshalJSON(r io.Reader, js interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, js)
}

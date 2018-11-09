package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
)

// ReadUnmarshalJSON 读取数据反序列化成JSON对象
func ReadUnmarshalJSON(r io.Reader, js interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, js)
}

// HTTPGetJSON GET请求返回数据反序列化成JSON对象
func HTTPGetJSON(url string, js interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return ReadUnmarshalJSON(resp.Body, js)
}

// RoundPrice 价格取整
func RoundPrice(v float64) float64 {
	return math.Round(v*100) / 100
}

// Bargain 随机砍价
func Bargain(m, n float64) float64 {
	if n == 1 {
		return m
	}
	max := m / n * 2
	res := rand.Float64() * max
	if res < 0.02 {
		return 0.01
	}
	return math.Floor(res*100) / 100
}

package xnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type xHttp struct{}

var XHttp = &xHttp{}

// 流式
func (x *xHttp) PostFormStreamChan(ctx context.Context, url string, reqData interface{}, HeaderValue http.Header) (chan string, error) {
	cli := &http.Client{
		Timeout: time.Second * 60,
	}

	var reqBody io.Reader
	if reqData != nil {
		jsonData, err := json.Marshal(reqData)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	if len(HeaderValue) > 0 {
		for k, v := range HeaderValue {
			req.Header.Set(k, v[0])
		}
	}

	var lines = make(chan string, 4)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			resp.Body.Close()
			close(lines)
		}()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Text()
				if line == "" {
					continue
				}
				lines <- line
			}

		}
	}()

	return lines, nil
}

type proxy struct {
	client *http.Client
}

var Proxy = &proxy{}

const LocalProxyUrl = "http://127.0.0.1:7890"

func SetProxyUrl(u string, timeout ...time.Duration) error {
	u = strings.TrimSpace(u)
	var t = 10 * time.Second
	if len(timeout) > 0 && timeout[0] > 0 {
		t = timeout[0]
	}

	var tans *http.Transport
	if u != "" {
		roxyURL, err := url.Parse(u) // 替
		if err != nil {
			return err
		}

		tans = &http.Transport{Proxy: http.ProxyURL(roxyURL)}
	}

	client := &http.Client{
		Transport: tans,
		Timeout:   t,
	}
	Proxy.client = client
	return nil
}

//func init() {
//	roxyURL, err := url.Parse("http://127.0.0.1:7890") // 替换为你的代理地址
//	if err != nil {
//		log.Fatalf("代理地址解析失败: %v", err)
//	}
//
//	// 创建带有代理的 HTTP 客户端
//	client := &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(roxyURL),
//		},
//		Timeout: 10 * time.Second,
//	}
//	Proxy.client = client
//}

func (p *proxy) Get(url string) (resp *http.Response, err error) {
	return p.client.Get(url)
}

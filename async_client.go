package primp

import (
	"sync"
)

// AsyncClient 是异步执行请求的客户端
type AsyncClient struct {
	*Client
}

// NewAsyncClient 创建新的 AsyncClient
func NewAsyncClient(options ...Option) *AsyncClient {
	return &AsyncClient{
		Client: NewClient(options...),
	}
}

// AsyncResponse 表示异步 HTTP 响应
type AsyncResponse struct {
	Response *Response
	Error    error
}

// RequestAsync 异步执行 HTTP 请求
func (c *AsyncClient) RequestAsync(method HttpMethod, url string, params RequestParams) <-chan AsyncResponse {
	ch := make(chan AsyncResponse, 1)

	go func() {
		resp, err := c.Client.Request(method, url, params)
		ch <- AsyncResponse{
			Response: resp,
			Error:    err,
		}
		close(ch)
	}()

	return ch
}

// GetAsync 异步发送 GET 请求
func (c *AsyncClient) GetAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(GET, url, reqParams)
}

// HeadAsync 异步发送 HEAD 请求
func (c *AsyncClient) HeadAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(HEAD, url, reqParams)
}

// OptionsAsync 异步发送 OPTIONS 请求
func (c *AsyncClient) OptionsAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(OPTIONS, url, reqParams)
}

// DeleteAsync 异步发送 DELETE 请求
func (c *AsyncClient) DeleteAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(DELETE, url, reqParams)
}

// PostAsync 异步发送 POST 请求
func (c *AsyncClient) PostAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(POST, url, reqParams)
}

// PutAsync 异步发送 PUT 请求
func (c *AsyncClient) PutAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(PUT, url, reqParams)
}

// PatchAsync 异步发送 PATCH 请求
func (c *AsyncClient) PatchAsync(url string, params ...RequestParams) <-chan AsyncResponse {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.RequestAsync(PATCH, url, reqParams)
}

// Batch 表示要并发执行的请求批次
type Batch struct {
	client    *AsyncClient
	wg        sync.WaitGroup
	mutex     sync.Mutex
	responses map[string]AsyncResponse
}

// NewBatch 创建新的请求批次
func (c *AsyncClient) NewBatch() *Batch {
	return &Batch{
		client:    c,
		responses: make(map[string]AsyncResponse),
	}
}

// Add 向批次添加请求
func (b *Batch) Add(id string, method HttpMethod, url string, params RequestParams) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		resp, err := b.client.Client.Request(method, url, params)
		b.mutex.Lock()
		defer b.mutex.Unlock()
		b.responses[id] = AsyncResponse{
			Response: resp,
			Error:    err,
		}
	}()
}

// Wait 等待批次中的所有请求完成
func (b *Batch) Wait() map[string]AsyncResponse {
	b.wg.Wait()
	return b.responses
}

package primp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Client 表示可以模拟各种浏览器的 HTTP 客户端
type Client struct {
	httpClient    *http.Client
	headers       map[string]string
	auth          *BasicAuth
	authBearer    string
	params        map[string]string
	proxy         string
	timeout       time.Duration
	impersonate   Impersonate
	impersonateOS ImpersonateOS
	cookieStore   bool
	referer       bool
	verify        bool
	caCertFile    string
	httpsOnly     bool
	http2Only     bool
}

// NewClient 创建一个新的带有给定选项的 HTTP 客户端
func NewClient(options ...Option) *Client {
	// 创建 cookie jar
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	// 默认客户端
	client := &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		headers:     make(map[string]string),
		cookieStore: true,
		referer:     true,
		verify:      true,
		timeout:     30 * time.Second,
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	// 应用浏览器模拟
	if client.impersonate != "" {
		client.applyBrowserImpersonation()
	}

	// 设置代理
	if client.proxy == "" {
		client.proxy = os.Getenv("PRIMP_PROXY")
	}
	if client.proxy != "" {
		proxyURL, err := url.Parse(client.proxy)
		if err == nil {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client.httpClient.Transport = transport
		}
	}

	// 应用 CA 证书（如果 verify 为 true）
	if client.verify {
		if client.caCertFile == "" {
			client.caCertFile = os.Getenv("PRIMP_CA_BUNDLE")
			if client.caCertFile == "" {
				client.caCertFile = os.Getenv("CA_CERT_FILE")
			}
		}
		if client.caCertFile != "" {
			client.applyCACertificate()
		}
	} else {
		// 如果 verify 为 false，则跳过 SSL 验证
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.httpClient.Transport = transport
	}

	return client
}

// applyBrowserImpersonation 设置用于模拟指定浏览器的头部
func (c *Client) applyBrowserImpersonation() {
	headers := getBrowserHeaders(c.impersonate, c.impersonateOS)
	for k, v := range headers {
		c.headers[k] = v
	}
}

// applyCACertificate 加载并应用 CA 证书
func (c *Client) applyCACertificate() {
	certPool, err := LoadCACerts(c.caCertFile)
	if err == nil && certPool != nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		}
		c.httpClient.Transport = transport
	}
}

// SetHeaders 设置请求头
func (c *Client) SetHeaders(headers map[string]string) {
	c.headers = headers
}

// Headers 返回当前头部
func (c *Client) Headers() map[string]string {
	return c.headers
}

// GetCookies 返回给定 URL 的 cookies
func (c *Client) GetCookies(urlStr string) (map[string]string, error) {
	if !c.cookieStore {
		return nil, nil
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	cookies := c.httpClient.Jar.Cookies(parsedURL)
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	return cookieMap, nil
}

// SetCookies 设置给定 URL 的 cookies
func (c *Client) SetCookies(urlStr string, cookies map[string]string) error {
	if !c.cookieStore {
		return nil
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	var httpCookies []*http.Cookie
	for name, value := range cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:  name,
			Value: value,
			Path:  "/",
		})
	}

	c.httpClient.Jar.SetCookies(parsedURL, httpCookies)
	return nil
}

// SetProxy 设置代理 URL
func (c *Client) SetProxy(proxyURL string) error {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	c.httpClient.Transport = transport
	c.proxy = proxyURL
	return nil
}

// SetImpersonate 设置要模拟的浏览器
func (c *Client) SetImpersonate(impersonate Impersonate) {
	c.impersonate = impersonate
	c.applyBrowserImpersonation()
}

// Impersonate 返回当前浏览器模拟设置
func (c *Client) GetImpersonate() Impersonate {
	return c.impersonate
}

// SetImpersonateOS 设置要模拟的操作系统
func (c *Client) SetImpersonateOS(impersonateOS ImpersonateOS) {
	c.impersonateOS = impersonateOS
	c.applyBrowserImpersonation()
}

// ImpersonateOS 返回当前操作系统模拟设置
func (c *Client) GetImpersonateOS() ImpersonateOS {
	return c.impersonateOS
}

// Request 使用指定方法和 URL 执行 HTTP 请求
func (c *Client) Request(method HttpMethod, urlStr string, params RequestParams) (*Response, error) {
	// 创建带有超时的上下文
	ctx := context.Background()
	if params.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, params.Timeout)
		defer cancel()
	} else if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 准备带有查询参数的 URL
	reqURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// 添加查询参数
	q := reqURL.Query()
	if params.Params != nil {
		for k, v := range params.Params {
			q.Add(k, v)
		}
	}
	if c.params != nil {
		for k, v := range c.params {
			q.Add(k, v)
		}
	}
	reqURL.RawQuery = q.Encode()

	// 创建请求体
	var body io.Reader
	var contentType string

	// 仅当方法为 POST、PUT 或 PATCH 时
	if method == POST || method == PUT || method == PATCH {
		if params.Content != nil {
			body = bytes.NewReader(params.Content)
		} else if params.Data != nil {
			formData := url.Values{}
			for k, v := range params.Data {
				formData.Add(k, fmt.Sprintf("%v", v))
			}
			body = strings.NewReader(formData.Encode())
			contentType = "application/x-www-form-urlencoded"
		} else if params.JSON != nil {
			jsonData, err := json.Marshal(params.JSON)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON: %w", err)
			}
			body = bytes.NewReader(jsonData)
			contentType = "application/json"
		} else if params.Files != nil {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			for fieldName, filePath := range params.Files {
				file, err := os.Open(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
				}
				defer file.Close()

				part, err := w.CreateFormFile(fieldName, filepath.Base(filePath))
				if err != nil {
					return nil, fmt.Errorf("failed to create form file: %w", err)
				}

				_, err = io.Copy(part, file)
				if err != nil {
					return nil, fmt.Errorf("failed to copy file content: %w", err)
				}
			}

			err = w.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to close multipart writer: %w", err)
			}

			body = &b
			contentType = w.FormDataContentType()
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, string(method), reqURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置内容类型
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置头部
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Set(k, v)
		}
	}

	// 设置 cookies
	if params.Cookies != nil {
		var cookieStrings []string
		for k, v := range params.Cookies {
			cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", k, v))
		}
		req.Header.Set("Cookie", strings.Join(cookieStrings, "; "))
	}

	// 设置认证
	if params.Auth != nil {
		req.SetBasicAuth(params.Auth.Username, params.Auth.Password)
	} else if c.auth != nil {
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	}

	// 设置 Bearer 令牌
	if params.AuthBearer != "" {
		req.Header.Set("Authorization", "Bearer "+params.AuthBearer)
	} else if c.authBearer != "" {
		req.Header.Set("Authorization", "Bearer "+c.authBearer)
	}

	// 设置 referer（如果启用）
	if c.referer && req.Header.Get("Referer") == "" && req.URL.Path != "/" {
		referer := fmt.Sprintf("%s://%s/", req.URL.Scheme, req.URL.Host)
		req.Header.Set("Referer", referer)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// 创建响应
	return newResponse(resp, reqURL.String())
}

// Get 发送 GET 请求
func (c *Client) Get(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(GET, url, reqParams)
}

// Head 发送 HEAD 请求
func (c *Client) Head(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(HEAD, url, reqParams)
}

// Options 发送 OPTIONS 请求
func (c *Client) Options(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(OPTIONS, url, reqParams)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(DELETE, url, reqParams)
}

// Post 发送 POST 请求
func (c *Client) Post(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(POST, url, reqParams)
}

// Put 发送 PUT 请求
func (c *Client) Put(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(PUT, url, reqParams)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string, params ...RequestParams) (*Response, error) {
	var reqParams RequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}
	return c.Request(PATCH, url, reqParams)
}

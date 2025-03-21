// Package primp 是一个能够模拟各种浏览器的 HTTP 客户端
package primp

import (
	"time"
)

// HttpMethod 表示 HTTP 请求方法
type HttpMethod string

const (
	GET     HttpMethod = "GET"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
	DELETE  HttpMethod = "DELETE"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	PATCH   HttpMethod = "PATCH"
)

// Impersonate 表示可以模拟的浏览器类型
type Impersonate string

// ImpersonateOS 表示可以模拟的操作系统
type ImpersonateOS string

// 支持的浏览器模拟类型常量定义
const (
	Chrome100 Impersonate = "chrome_100"
	Chrome101 Impersonate = "chrome_101"
	Chrome104 Impersonate = "chrome_104"
	Chrome105 Impersonate = "chrome_105"
	Chrome106 Impersonate = "chrome_106"
	Chrome107 Impersonate = "chrome_107"
	Chrome108 Impersonate = "chrome_108"
	Chrome109 Impersonate = "chrome_109"
	Chrome114 Impersonate = "chrome_114"
	Chrome116 Impersonate = "chrome_116"
	Chrome117 Impersonate = "chrome_117"
	Chrome118 Impersonate = "chrome_118"
	Chrome119 Impersonate = "chrome_119"
	Chrome120 Impersonate = "chrome_120"
	Chrome123 Impersonate = "chrome_123"
	Chrome124 Impersonate = "chrome_124"
	Chrome126 Impersonate = "chrome_126"
	Chrome127 Impersonate = "chrome_127"
	Chrome128 Impersonate = "chrome_128"
	Chrome129 Impersonate = "chrome_129"
	Chrome130 Impersonate = "chrome_130"
	Chrome131 Impersonate = "chrome_131"
	Chrome133 Impersonate = "chrome_133"

	SafariIos165  Impersonate = "safari_ios_16.5"
	SafariIos172  Impersonate = "safari_ios_17.2"
	SafariIos1741 Impersonate = "safari_ios_17.4.1"
	SafariIos1811 Impersonate = "safari_ios_18.1.1"
	SafariIPad18  Impersonate = "safari_ipad_18"
	Safari153     Impersonate = "safari_15.3"
	Safari155     Impersonate = "safari_15.5"
	Safari1561    Impersonate = "safari_15.6.1"
	Safari16      Impersonate = "safari_16"
	Safari165     Impersonate = "safari_16.5"
	Safari170     Impersonate = "safari_17.0"
	Safari1721    Impersonate = "safari_17.2.1"
	Safari1741    Impersonate = "safari_17.4.1"
	Safari175     Impersonate = "safari_17.5"
	Safari18      Impersonate = "safari_18"
	Safari182     Impersonate = "safari_18.2"

	OkHttp39  Impersonate = "okhttp_3.9"
	OkHttp311 Impersonate = "okhttp_3.11"
	OkHttp313 Impersonate = "okhttp_3.13"
	OkHttp314 Impersonate = "okhttp_3.14"
	OkHttp49  Impersonate = "okhttp_4.9"
	OkHttp410 Impersonate = "okhttp_4.10"
	OkHttp5   Impersonate = "okhttp_5"

	Edge101 Impersonate = "edge_101"
	Edge122 Impersonate = "edge_122"
	Edge127 Impersonate = "edge_127"
	Edge131 Impersonate = "edge_131"

	Firefox109 Impersonate = "firefox_109"
	Firefox117 Impersonate = "firefox_117"
	Firefox128 Impersonate = "firefox_128"
	Firefox133 Impersonate = "firefox_133"
	Firefox135 Impersonate = "firefox_135"
)

// 支持的操作系统模拟类型常量定义
const (
	Android ImpersonateOS = "android"
	IOS     ImpersonateOS = "ios"
	Linux   ImpersonateOS = "linux"
	MacOS   ImpersonateOS = "macos"
	Windows ImpersonateOS = "windows"
)

// RequestParams 定义请求可用的参数
type RequestParams struct {
	Auth       *BasicAuth
	AuthBearer string
	Params     map[string]string
	Headers    map[string]string
	Cookies    map[string]string
	Timeout    time.Duration
	Content    []byte
	Data       map[string]interface{}
	JSON       interface{}
	Files      map[string]string
}

// ClientRequestParams 扩展 RequestParams 添加客户端特定选项
type ClientRequestParams struct {
	RequestParams
	Impersonate   Impersonate
	ImpersonateOS ImpersonateOS
	Verify        bool
	CACertFile    string
}

// BasicAuth 表示 HTTP 基本认证凭据
type BasicAuth struct {
	Username string
	Password string
}

// 便捷的 HTTP 方法函数
func Get(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(GET, url, params...)
}

func Head(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(HEAD, url, params...)
}

func Options(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(OPTIONS, url, params...)
}

func Delete(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(DELETE, url, params...)
}

func Post(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(POST, url, params...)
}

func Put(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(PUT, url, params...)
}

func Patch(url string, params ...ClientRequestParams) (*Response, error) {
	return Request(PATCH, url, params...)
}

// Request 使用指定方法和 URL 执行 HTTP 请求
func Request(method HttpMethod, url string, params ...ClientRequestParams) (*Response, error) {
	var reqParams ClientRequestParams
	if len(params) > 0 {
		reqParams = params[0]
	}

	client := NewClient(
		WithImpersonate(reqParams.Impersonate),
		WithImpersonateOS(reqParams.ImpersonateOS),
		WithVerify(reqParams.Verify),
		WithCACertFile(reqParams.CACertFile),
	)

	return client.Request(method, url, reqParams.RequestParams)
}

// Option 是配置 Client 的函数类型
type Option func(*Client)

// WithImpersonate 设置要模拟的浏览器
func WithImpersonate(impersonate Impersonate) Option {
	return func(c *Client) {
		c.impersonate = impersonate
	}
}

// WithImpersonateOS 设置要模拟的操作系统
func WithImpersonateOS(impersonateOS ImpersonateOS) Option {
	return func(c *Client) {
		c.impersonateOS = impersonateOS
	}
}

// WithVerify 启用或禁用 SSL 验证
func WithVerify(verify bool) Option {
	return func(c *Client) {
		c.verify = verify
	}
}

// WithCACertFile 设置自定义 CA 证书文件
func WithCACertFile(caCertFile string) Option {
	return func(c *Client) {
		c.caCertFile = caCertFile
	}
}

package primp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Response 表示 HTTP 响应
type Response struct {
	httpResp   *http.Response
	content    []byte
	encoding   string
	headers    map[string]string
	cookies    map[string]string
	URL        string
	StatusCode int
}

// newResponse 从 http.Response 创建新的 Response
func newResponse(resp *http.Response, url string) (*Response, error) {
	return &Response{
		httpResp:   resp,
		URL:        url,
		StatusCode: resp.StatusCode,
	}, nil
}

// Content 以字节形式返回响应体
func (r *Response) Content() ([]byte, error) {
	if r.content != nil {
		return r.content, nil
	}

	defer r.httpResp.Body.Close()
	content, err := io.ReadAll(r.httpResp.Body)
	if err != nil {
		return nil, err
	}

	r.content = content
	return content, nil
}

// Text 以字符串形式返回响应体
func (r *Response) Text() (string, error) {
	content, err := r.Content()
	if err != nil {
		return "", err
	}

	encoding := r.Encoding()
	dec, err := getDecoder(encoding)
	if err != nil {
		// 默认使用 UTF-8
		return string(content), nil
	}

	reader := transform.NewReader(bytes.NewReader(content), dec.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// getDecoder 返回给定字符集的编码
func getDecoder(charset string) (encoding.Encoding, error) {
	if charset == "" || strings.EqualFold(charset, "utf-8") || strings.EqualFold(charset, "utf8") {
		return unicode.UTF8, nil
	}
	if strings.EqualFold(charset, "latin1") || strings.EqualFold(charset, "iso-8859-1") {
		return charmap.ISO8859_1, nil
	}
	return htmlindex.Get(charset)
}

// JSON 将响应体反序列化到提供的值中
func (r *Response) JSON(v interface{}) error {
	content, err := r.Content()
	if err != nil {
		return err
	}

	return json.Unmarshal(content, v)
}

// Headers 返回响应头
func (r *Response) Headers() (map[string]string, error) {
	if r.headers != nil {
		return r.headers, nil
	}

	headers := make(map[string]string)
	for name, values := range r.httpResp.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	r.headers = headers
	return headers, nil
}

// Cookies 返回响应 cookies
func (r *Response) Cookies() (map[string]string, error) {
	if r.cookies != nil {
		return r.cookies, nil
	}

	cookies := make(map[string]string)
	for _, cookie := range r.httpResp.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	r.cookies = cookies
	return cookies, nil
}

// Encoding 返回响应编码
func (r *Response) Encoding() string {
	if r.encoding != "" {
		return r.encoding
	}

	contentType := r.httpResp.Header.Get("Content-Type")
	if contentType != "" {
		parts := strings.Split(contentType, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "charset=") {
				r.encoding = strings.TrimPrefix(part, "charset=")
				r.encoding = strings.Trim(r.encoding, `"'`)
				return r.encoding
			}
		}
	}

	// 默认为 UTF-8
	r.encoding = "utf-8"
	return r.encoding
}

// SetEncoding 设置响应编码
func (r *Response) SetEncoding(encoding string) {
	r.encoding = encoding
}

// Stream 返回响应体的读取器
func (r *Response) Stream() (io.ReadCloser, error) {
	return r.httpResp.Body, nil
}

// TextMarkdown 以 Markdown 文本形式返回响应
func (r *Response) TextMarkdown() (string, error) {
	content, err := r.Content()
	if err != nil {
		return "", err
	}

	// 简单的 HTML 到 Markdown 转换
	return htmlToMarkdown(content), nil
}

// TextPlain 以纯文本形式返回响应
func (r *Response) TextPlain() (string, error) {
	content, err := r.Content()
	if err != nil {
		return "", err
	}

	// 简单的 HTML 到纯文本转换
	return htmlToPlainText(content), nil
}

// TextRich 以富文本形式返回响应
func (r *Response) TextRich() (string, error) {
	content, err := r.Content()
	if err != nil {
		return "", err
	}

	// 简单的 HTML 到富文本转换
	return htmlToRichText(content), nil
}

// htmlToMarkdown 将 HTML 转换为 Markdown
func htmlToMarkdown(htmlContent []byte) string {
	// 这是一个简化的实现
	// 在实际实现中，应该使用适当的 HTML 到 Markdown 转换器
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return string(htmlContent)
	}

	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return sb.String()
}

// htmlToPlainText 将 HTML 转换为纯文本
func htmlToPlainText(htmlContent []byte) string {
	// 这是一个简化的实现
	// 在实际实现中，应该使用适当的 HTML 到文本转换器
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return string(htmlContent)
	}

	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return sb.String()
}

// htmlToRichText 将 HTML 转换为富文本
func htmlToRichText(htmlContent []byte) string {
	// 这是一个简化的实现
	// 在实际实现中，应该使用适当的 HTML 到富文本转换器
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return string(htmlContent)
	}

	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return sb.String()
}

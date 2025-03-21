package primp

import (
	"fmt"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
)

// Impersonate和ImpersonateOS类型保持不变

// 浏览器配置，包含所有模拟参数
type BrowserProfile struct {
	UserAgent string
	Headers   map[string]string
}

// 获取浏览器配置的主函数
func GetBrowserProfile(browserName string, osName string) (BrowserProfile, error) {
	browser, err := ImpersonateFromString(browserName)
	if err != nil {
		return BrowserProfile{}, err
	}

	os, err := ImpersonateOSFromString(osName)
	if err != nil {
		return BrowserProfile{}, err
	}

	headers := getBrowserHeaders(browser, os)
	userAgent := headers["User-Agent"]

	return BrowserProfile{
		UserAgent: userAgent,
		Headers:   headers,
	}, nil
}

// 获取浏览器头信息
func getBrowserHeaders(browser Impersonate, os ImpersonateOS) map[string]string {
	// 如果没有指定操作系统，默认为Windows
	if os == "" {
		os = Windows
	}

	// 获取基本的通用头部
	headers := getBaseHeaders(browser)

	// 添加User-Agent
	headers["User-Agent"] = getUserAgent(browser, os)

	// 添加特定于浏览器的附加头部
	switch {
	case strings.HasPrefix(string(browser), "Chrome"):
		addChromeHeaders(headers, browser, os)
	case strings.HasPrefix(string(browser), "Firefox"):
		addFirefoxHeaders(headers)
	case strings.HasPrefix(string(browser), "Safari"):
		addSafariHeaders(headers)
	}

	return headers
}

// 获取基本头信息
func getBaseHeaders(browser Impersonate) map[string]string {
	return map[string]string{
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
	}
}

// 添加Chrome特有的头信息
func addChromeHeaders(headers map[string]string, browser Impersonate, os ImpersonateOS) {
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"

	// 从浏览器字符串中提取版本号
	version := extractVersion(string(browser))

	headers["sec-ch-ua"] = fmt.Sprintf(`"Google Chrome";v="%s", "Chromium";v="%s", "Not-A.Brand";v="99"`, version, version)
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-ch-ua-platform"] = `"` + string(os) + `"`
	headers["Sec-Fetch-Dest"] = "document"
	headers["Sec-Fetch-Mode"] = "navigate"
	headers["Sec-Fetch-Site"] = "none"
	headers["Sec-Fetch-User"] = "?1"
	headers["Upgrade-Insecure-Requests"] = "1"
}

// 添加Firefox特有的头信息
func addFirefoxHeaders(headers map[string]string) {
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	headers["Accept-Language"] = "en-US,en;q=0.5"
	headers["Upgrade-Insecure-Requests"] = "1"
	headers["Sec-Fetch-Dest"] = "document"
	headers["Sec-Fetch-Mode"] = "navigate"
	headers["Sec-Fetch-Site"] = "none"
	headers["Sec-Fetch-User"] = "?1"
}

// 添加Safari特有的头信息
func addSafariHeaders(headers map[string]string) {
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	headers["Accept-Language"] = "en-US,en;q=0.9"
}

// 从浏览器名称中提取版本号
func extractVersion(browserString string) string {
	parts := strings.Split(browserString, "_")
	if len(parts) > 1 {
		return parts[1]
	}
	return "133" // 默认版本
}

// 获取用户代理字符串
func getUserAgent(browserType Impersonate, os ImpersonateOS) string {
	browserName := strings.Split(string(browserType), "_")[0]

	var baseUA string
	switch strings.ToLower(browserName) {
	case "chrome":
		baseUA = browser.Chrome()
	case "firefox":
		baseUA = browser.Firefox()
	case "safari":
		baseUA = browser.Safari()
	default:
		baseUA = browser.Chrome()
	}

	// 修改用户代理的版本信息以匹配指定版本
	version := extractVersion(string(browserType))
	ua := modifyUserAgentVersion(baseUA, browserName, version, string(os))

	return ua
}

// 修改用户代理版本信息
func modifyUserAgentVersion(baseUA, browserName, version, os string) string {
	switch strings.ToLower(browserName) {
	case "chrome":
		if strings.Contains(baseUA, "Chrome/") {
			parts := strings.Split(baseUA, "Chrome/")
			versionParts := strings.Split(parts[1], " ")
			return parts[0] + "Chrome/" + version + ".0.0.0 " + strings.Join(versionParts[1:], " ")
		}
	case "firefox":
		if strings.Contains(baseUA, "Firefox/") {
			parts := strings.Split(baseUA, "Firefox/")
			return parts[0] + "Firefox/" + version + ".0"
		}
	}

	return baseUA
}

// ImpersonateFromString 和 ImpersonateOSFromString 保持不变

// ImpersonateFromString 将字符串解析为 Impersonate 值
func ImpersonateFromString(s string) (Impersonate, error) {
	switch strings.ToLower(s) {
	case "chrome_100":
		return Chrome100, nil
	case "chrome_101":
		return Chrome101, nil
	case "chrome_104":
		return Chrome104, nil
	case "chrome_105":
		return Chrome105, nil
	case "chrome_106":
		return Chrome106, nil
	case "chrome_107":
		return Chrome107, nil
	case "chrome_108":
		return Chrome108, nil
	case "chrome_109":
		return Chrome109, nil
	case "chrome_114":
		return Chrome114, nil
	case "chrome_116":
		return Chrome116, nil
	case "chrome_117":
		return Chrome117, nil
	case "chrome_118":
		return Chrome118, nil
	case "chrome_119":
		return Chrome119, nil
	case "chrome_120":
		return Chrome120, nil
	case "chrome_123":
		return Chrome123, nil
	case "chrome_124":
		return Chrome124, nil
	case "chrome_126":
		return Chrome126, nil
	case "chrome_127":
		return Chrome127, nil
	case "chrome_128":
		return Chrome128, nil
	case "chrome_129":
		return Chrome129, nil
	case "chrome_130":
		return Chrome130, nil
	case "chrome_131":
		return Chrome131, nil
	case "chrome_133":
		return Chrome133, nil
	case "safari_ios_16.5":
		return SafariIos165, nil
	case "safari_ios_17.2":
		return SafariIos172, nil
	case "safari_ios_17.4.1":
		return SafariIos1741, nil
	case "safari_ios_18.1.1":
		return SafariIos1811, nil
	case "safari_ipad_18":
		return SafariIPad18, nil
	case "safari_15.3":
		return Safari153, nil
	case "safari_15.5":
		return Safari155, nil
	case "safari_15.6.1":
		return Safari1561, nil
	case "safari_16":
		return Safari16, nil
	case "safari_16.5":
		return Safari165, nil
	case "safari_17.0":
		return Safari170, nil
	case "safari_17.2.1":
		return Safari1721, nil
	case "safari_17.4.1":
		return Safari1741, nil
	case "safari_17.5":
		return Safari175, nil
	case "safari_18":
		return Safari18, nil
	case "safari_18.2":
		return Safari182, nil
	case "okhttp_3.9":
		return OkHttp39, nil
	case "okhttp_3.11":
		return OkHttp311, nil
	case "okhttp_3.13":
		return OkHttp313, nil
	case "okhttp_3.14":
		return OkHttp314, nil
	case "okhttp_4.9":
		return OkHttp49, nil
	case "okhttp_4.10":
		return OkHttp410, nil
	case "okhttp_5":
		return OkHttp5, nil
	case "firefox_109":
		return Firefox109, nil
	case "firefox_117":
		return Firefox117, nil
	case "firefox_128":
		return Firefox128, nil
	case "firefox_133":
		return Firefox133, nil
	case "firefox_135":
		return Firefox135, nil
	default:
		return "", fmt.Errorf("invalid impersonate: %s", s)
	}
}

// ImpersonateOSFromString 将字符串解析为 ImpersonateOS 值
func ImpersonateOSFromString(s string) (ImpersonateOS, error) {
	switch strings.ToLower(s) {
	case "android":
		return Android, nil
	case "ios":
		return IOS, nil
	case "linux":
		return Linux, nil
	case "macos":
		return MacOS, nil
	case "windows":
		return Windows, nil
	default:
		return "", fmt.Errorf("invalid impersonate_os: %s", s)
	}
}

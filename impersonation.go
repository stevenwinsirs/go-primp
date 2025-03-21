package primp

import (
	"fmt"
	"strings"
)

// getBrowserHeaders 返回用于模拟特定浏览器和操作系统的头部
func getBrowserHeaders(browser Impersonate, os ImpersonateOS) map[string]string {
	// 如果没有指定操作系统，默认为 Windows
	if os == "" {
		os = Windows
	}

	// 大多数浏览器通用的基本头部
	headers := map[string]string{
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
	}

	// 添加浏览器特定头部
	switch browser {
	case Chrome133:
		headers["User-Agent"] = getUserAgent(browser, os)
		headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
		headers["sec-ch-ua"] = `"Google Chrome";v="133", "Chromium";v="133", "Not-A.Brand";v="99"`
		headers["sec-ch-ua-mobile"] = "?0"
		headers["sec-ch-ua-platform"] = `"` + string(os) + `"`
		headers["Sec-Fetch-Dest"] = "document"
		headers["Sec-Fetch-Mode"] = "navigate"
		headers["Sec-Fetch-Site"] = "none"
		headers["Sec-Fetch-User"] = "?1"
		headers["Upgrade-Insecure-Requests"] = "1"
	case Chrome131:
		headers["User-Agent"] = getUserAgent(browser, os)
		headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
		headers["sec-ch-ua"] = `"Google Chrome";v="131", "Chromium";v="131", "Not-A.Brand";v="99"`
		headers["sec-ch-ua-mobile"] = "?0"
		headers["sec-ch-ua-platform"] = `"` + string(os) + `"`
		headers["Sec-Fetch-Dest"] = "document"
		headers["Sec-Fetch-Mode"] = "navigate"
		headers["Sec-Fetch-Site"] = "none"
		headers["Sec-Fetch-User"] = "?1"
		headers["Upgrade-Insecure-Requests"] = "1"
	// ... (添加更多浏览器实现)
	case Firefox135:
		headers["User-Agent"] = getUserAgent(browser, os)
		headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
		headers["Accept-Language"] = "en-US,en;q=0.5"
		headers["Connection"] = "keep-alive"
		headers["Upgrade-Insecure-Requests"] = "1"
		headers["Sec-Fetch-Dest"] = "document"
		headers["Sec-Fetch-Mode"] = "navigate"
		headers["Sec-Fetch-Site"] = "none"
		headers["Sec-Fetch-User"] = "?1"
	case Safari18:
		headers["User-Agent"] = getUserAgent(browser, os)
		headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
		headers["Accept-Language"] = "en-US,en;q=0.9"
		headers["Connection"] = "keep-alive"
		// ... (添加更多浏览器实现)
	}

	return headers
}

// getUserAgent 返回浏览器和操作系统的适当 User-Agent 字符串
func getUserAgent(browser Impersonate, os ImpersonateOS) string {
	switch browser {
	case Chrome133:
		switch os {
		case Windows:
			return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
		case MacOS:
			return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
		case Linux:
			return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
		case Android:
			return "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36"
		case IOS:
			return "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/133.0.0.0 Mobile/15E148 Safari/604.1"
		}
	case Chrome131:
		switch os {
		case Windows:
			return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
		case MacOS:
			return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
		case Linux:
			return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
		case Android:
			return "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36"
		case IOS:
			return "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/131.0.0.0 Mobile/15E148 Safari/604.1"
		}
	// ... (添加更多浏览器 User-Agent 实现)
	case Firefox135:
		switch os {
		case Windows:
			return "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0"
		case MacOS:
			return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0"
		case Linux:
			return "Mozilla/5.0 (X11; Linux x86_64; rv:135.0) Gecko/20100101 Firefox/135.0"
		case Android:
			return "Mozilla/5.0 (Android 10; Mobile; rv:135.0) Gecko/135.0 Firefox/135.0"
		case IOS:
			return "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/135.0 Mobile/15E148 Safari/605.1.15"
		}
	case Safari18:
		switch os {
		case MacOS:
			return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15"
		case IOS:
			return "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1"
		default:
			// Safari 主要在 macOS 和 iOS 上
			return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15"
		}
	}

	// 默认返回最新 Chrome 版本的 Windows User-Agent
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
}

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
	case "edge_101":
		return Edge101, nil
	case "edge_122":
		return Edge122, nil
	case "edge_127":
		return Edge127, nil
	case "edge_131":
		return Edge131, nil
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

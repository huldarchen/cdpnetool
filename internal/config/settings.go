package config

// DefaultSettings 定义所有设置的默认值
type DefaultSettings struct {
	Language    string
	Theme       string
	DevToolsURL string
	BrowserArgs string
	BrowserPath string
}

// GetDefaultSettings 返回默认设置
func GetDefaultSettings() DefaultSettings {
	return DefaultSettings{
		Language:    "zh",
		Theme:       "system",
		DevToolsURL: "http://localhost:9222",
		BrowserArgs: "", // 空表示无额外参数
		BrowserPath: "", // 空表示自动检测系统浏览器
	}
}

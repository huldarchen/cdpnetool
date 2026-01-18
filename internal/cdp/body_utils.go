package cdp

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	"github.com/mafredri/cdp/protocol/fetch"
)

// GetRequestBody 统一提取并解码请求体
func GetRequestBody(ev *fetch.RequestPausedReply) string {
	if ev == nil {
		return ""
	}

	// 1. 如果有 postData 直接返回（CDP 已处理为普通字符串）
	if ev.Request.PostData != nil {
		return *ev.Request.PostData
	}

	// 2. 如果有 postDataEntries，需要按条目解码并拼接
	if len(ev.Request.PostDataEntries) > 0 {
		var sb strings.Builder
		for _, entry := range ev.Request.PostDataEntries {
			if entry.Bytes != nil {
				// 关键：对每一段进行 Base64 解码
				decoded, err := base64.StdEncoding.DecodeString(*entry.Bytes)
				if err == nil {
					sb.Write(decoded)
				} else {
					// 解码失败则保留原始 Base64（兜底）
					sb.WriteString(*entry.Bytes)
				}
			}
		}
		return sb.String()
	}

	return ""
}

// IsTextualBody 判断 Body 是否为文本类型，以便安全展示或匹配
func IsTextualBody(data []byte, contentType string) bool {
	lc := strings.ToLower(contentType)
	// 常见的文本类型
	if strings.HasPrefix(lc, "text/") ||
		strings.Contains(lc, "json") ||
		strings.Contains(lc, "xml") ||
		strings.Contains(lc, "javascript") ||
		strings.Contains(lc, "x-www-form-urlencoded") {
		return true
	}

	// 启发式检测：如果是有效的 UTF-8 且不包含过多的控制字符
	return utf8.Valid(data)
}

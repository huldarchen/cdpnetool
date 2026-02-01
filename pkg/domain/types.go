package domain

import (
	"net/http"
	"strings"
)

// SessionID 会话ID
type SessionID string

// TargetID 目标ID
type TargetID string

// RuleID 规则ID
type RuleID string

// SessionConfig 会话配置
type SessionConfig struct {
	DevToolsURL       string `json:"devToolsURL"`
	Concurrency       int    `json:"concurrency"`
	BodySizeThreshold int64  `json:"bodySizeThreshold"`
	PendingCapacity   int    `json:"pendingCapacity"`
	ProcessTimeoutMS  int    `json:"processTimeoutMS"`
}

// EngineStats 引擎统计信息
type EngineStats struct {
	Total   int64            `json:"total"`
	Matched int64            `json:"matched"`
	ByRule  map[RuleID]int64 `json:"byRule"`
}

// TargetInfo 目标信息
type TargetInfo struct {
	ID        TargetID `json:"id"`
	Type      string   `json:"type"`
	URL       string   `json:"url"`
	Title     string   `json:"title"`
	IsCurrent bool     `json:"isCurrent"`
}

// Header 封装通用的头部操作
type Header map[string]string

// Get 获取指定 Header 的值（大小写不敏感）
func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	return h[strings.ToLower(key)]
}

// Set 设置指定 Header 的值（自动转换为小写）
func (h Header) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

// Del 删除指定 Header
func (h Header) Del(key string) {
	delete(h, strings.ToLower(key))
}

// Request 请求模型
type Request struct {
	ID           string            `json:"id"`                     // 事务唯一ID
	URL          string            `json:"url"`                    // 完整URL
	Method       string            `json:"method"`                 // HTTP方法
	Headers      Header            `json:"headers"`                // 请求头
	Body         []byte            `json:"body"`                   // 请求体原始数据
	ResourceType string            `json:"resourceType,omitempty"` // 资源类型 (如 Document, XHR)
	Query        map[string]string `json:"query,omitempty"`        // 预解析的查询参数
	Cookies      map[string]string `json:"cookies,omitempty"`      // 预解析的Cookie
}

// Response 响应模型
type Response struct {
	StatusCode int            `json:"statusCode"`
	Headers    Header         `json:"headers"`
	Body       []byte         `json:"body"`
	Timing     ResponseTiming `json:"timing,omitempty"`
}

// ResponseTiming 响应时间信息
type ResponseTiming struct {
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

// RuleMatch 规则匹配信息
type RuleMatch struct {
	RuleID   string   `json:"ruleId"`
	RuleName string   `json:"ruleName"`
	Actions  []string `json:"actions"`
}

// NetworkEvent 网络请求事件（统一所有拦截事件）
type NetworkEvent struct {
	ID           string      `json:"id"` // 事务唯一ID (CDP RequestID)
	Session      SessionID   `json:"session"`
	Target       TargetID    `json:"target"`
	Timestamp    int64       `json:"timestamp"`
	IsMatched    bool        `json:"isMatched"` // 是否匹配规则
	Request      Request     `json:"request"`
	Response     *Response   `json:"response,omitempty"`
	FinalResult  string      `json:"finalResult,omitempty"`  // blocked / modified / passed
	MatchedRules []RuleMatch `json:"matchedRules,omitempty"` // 匹配的规则列表
}

// NewRequest 创建初始化请求对象
func NewRequest() *Request {
	return &Request{
		Headers: make(Header),
		Query:   make(map[string]string),
		Cookies: make(map[string]string),
	}
}

// NewResponse 创建初始化响应对象
func NewResponse() *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Headers:    make(Header),
	}
}

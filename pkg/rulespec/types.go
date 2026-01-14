package rulespec

import "cdpnetool/pkg/model"

// ConditionType 条件类型
type ConditionType string

// ConditionMode 条件模式
type ConditionMode string

// ConditionOp 条件操作符
type ConditionOp string

// PauseStage 暂停阶段
type PauseStage string

// PauseDefaultActionType 暂停默认动作类型
type PauseDefaultActionType string

// RuleMode 规则模式
type RuleMode string

// JSONPatchOpType JSON补丁操作类型
type JSONPatchOpType string

const (
	JSONPatchOpAdd     JSONPatchOpType = "add"     // 添加
	JSONPatchOpRemove  JSONPatchOpType = "remove"  // 移除
	JSONPatchOpReplace JSONPatchOpType = "replace" // 替换
	JSONPatchOpMove    JSONPatchOpType = "move"    // 移动
	JSONPatchOpCopy    JSONPatchOpType = "copy"    // 复制
	JSONPatchOpTest    JSONPatchOpType = "test"    // 测试
)

const (
	ConditionTypeURL         ConditionType = "url"          // URL
	ConditionTypeMethod      ConditionType = "method"       // 方法
	ConditionTypeHeader      ConditionType = "header"       // 头部
	ConditionTypeQuery       ConditionType = "query"        // 查询
	ConditionTypeCookie      ConditionType = "cookie"       // Cookie
	ConditionTypeText        ConditionType = "text"         // 文本
	ConditionTypeMIME        ConditionType = "mime"         // MIME
	ConditionTypeSize        ConditionType = "size"         // 大小
	ConditionTypeProbability ConditionType = "probability"  // 概率
	ConditionTypeTimeWindow  ConditionType = "time_window"  // 时间窗口
	ConditionTypeJSONPointer ConditionType = "json_pointer" // JSON指针
	ConditionTypeStage       ConditionType = "stage"        // 阶段
)

const (
	ConditionModePrefix ConditionMode = "prefix" // 前缀
	ConditionModeRegex  ConditionMode = "regex"  // 正则
	ConditionModeExact  ConditionMode = "exact"  // 精确
)

const (
	ConditionOpEquals   ConditionOp = "equals"   // 等于
	ConditionOpContains ConditionOp = "contains" // 包含
	ConditionOpRegex    ConditionOp = "regex"    // 正则
	ConditionOpLT       ConditionOp = "lt"       // 小于
	ConditionOpLTE      ConditionOp = "lte"      // 小于等于
	ConditionOpGT       ConditionOp = "gt"       // 大于
	ConditionOpGTE      ConditionOp = "gte"      // 大于等于
	ConditionOpBetween  ConditionOp = "between"  // 区间
)

const (
	PauseStageRequest  PauseStage = "request"  // 请求
	PauseStageResponse PauseStage = "response" // 响应
)

const (
	PauseDefaultActionContinueOriginal PauseDefaultActionType = "continue_original" // 继续原始
	PauseDefaultActionContinueMutated  PauseDefaultActionType = "continue_mutated"  // 继续变异
	PauseDefaultActionFulfill          PauseDefaultActionType = "fulfill"           // 满足
	PauseDefaultActionFail             PauseDefaultActionType = "fail"              // 失败
)

const (
	RuleModeShortCircuit RuleMode = "short_circuit" // 短路
	RuleModeAggregate    RuleMode = "aggregate"     // 聚合
)

// RuleSet 规则集
type RuleSet struct {
	Version string `json:"version"`
	Rules   []Rule `json:"rules"`
}

// Rule 规则
type Rule struct {
	ID       model.RuleID `json:"id"`
	Name     string       `json:"name,omitempty"`
	Priority int          `json:"priority"`
	Mode     RuleMode     `json:"mode"`
	Match    Match        `json:"match"`
	Action   Action       `json:"action"`
}

// Match 匹配
type Match struct {
	AllOf  []Condition `json:"allOf"`
	AnyOf  []Condition `json:"anyOf"`
	NoneOf []Condition `json:"noneOf"`
}

// Condition 条件
type Condition struct {
	Type    ConditionType `json:"type"`
	Mode    ConditionMode `json:"mode"`
	Pattern string        `json:"pattern"`
	Values  []string      `json:"values"`
	Key     string        `json:"key"`
	Op      ConditionOp   `json:"op"`
	Value   string        `json:"value"`
	Pointer string        `json:"pointer"`
}

// Action 动作
type Action struct {
	Rewrite  *Rewrite `json:"rewrite"`
	Respond  *Respond `json:"respond"`
	Fail     *Fail    `json:"fail"`
	DelayMS  int      `json:"delayMS"`
	DropRate float64  `json:"dropRate"`
	Pause    *Pause   `json:"pause"`
}

// Rewrite 重写
type Rewrite struct {
	URL     *string            `json:"url"`
	Method  *string            `json:"method"`
	Headers map[string]*string `json:"headers"`
	Query   map[string]*string `json:"query"`
	Cookies map[string]*string `json:"cookies"`
	Body    *BodyPatch         `json:"body"`
}

// BodyPatch 体补丁
type BodyPatch struct {
	JSONPatch []JSONPatchOp   `json:"jsonPatch,omitempty"`
	TextRegex *TextRegexPatch `json:"textRegex,omitempty"`
	Base64    *Base64Patch    `json:"base64,omitempty"`
}

// JSONPatchOp JSON补丁操作
type JSONPatchOp struct {
	Op    JSONPatchOpType `json:"op"`
	Path  string          `json:"path"`
	From  string          `json:"from,omitempty"`
	Value any             `json:"value,omitempty"`
}

// TextRegexPatch 文本正则补丁
type TextRegexPatch struct {
	Pattern string `json:"pattern"`
	Replace string `json:"replace"`
}

// Base64Patch Base64补丁
type Base64Patch struct {
	Value string `json:"value"`
}

// Respond 响应
type Respond struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	Base64  bool              `json:"base64"`
}

// Fail 失败
type Fail struct {
	Reason string `json:"reason"`
}

// Pause 暂停
type Pause struct {
	Stage         PauseStage `json:"stage"`
	TimeoutMS     int        `json:"timeoutMS"`
	DefaultAction struct {
		Type   PauseDefaultActionType `json:"type"`
		Status int                    `json:"status"`
		Reason string                 `json:"reason"`
	} `json:"defaultAction"`
}

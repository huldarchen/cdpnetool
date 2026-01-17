package storage

import (
	"time"
)

// Setting 用户设置表
type Setting struct {
	Key       string    `gorm:"primaryKey" json:"key"`  // 设置键
	Value     string    `gorm:"type:text" json:"value"` // 设置值
	UpdatedAt time.Time `json:"updatedAt"`              // 更新时间
}

// 预定义的设置 Key
const (
	SettingKeyDevToolsURL  = "devtools_url"   // 开发者工具URL
	SettingKeyTheme        = "theme"          // 主题
	SettingKeyWindowBounds = "window_bounds"  // 窗口大小和位置
	SettingKeyLastConfigID = "last_config_id" // 上次使用的配置 ID
)

// ConfigRecord 配置表（存储规则配置）
type ConfigRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`             // 主键ID
	Name        string    `gorm:"uniqueIndex;not null" json:"name"` // 配置名称
	Description string    `gorm:"type:text" json:"description"`     // 配置描述
	Version     string    `json:"version"`                          // 配置格式版本
	RulesJSON   string    `gorm:"type:text" json:"rulesJson"`       // JSON 序列化的规则数组
	IsActive    bool      `gorm:"default:false" json:"isActive"`    // 是否为激活配置
	CreatedAt   time.Time `json:"createdAt"`                        // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`                        // 更新时间
}

// InterceptEventRecord 拦截事件历史表
type InterceptEventRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`   // 主键ID
	SessionID  string    `gorm:"index" json:"sessionId"` // 会话ID
	TargetID   string    `json:"targetId"`               // 目标ID
	Type       string    `gorm:"index" json:"type"`      // matched, rewritten, failed, rejected...
	URL        string    `json:"url"`                    // URL
	Method     string    `json:"method"`                 // 方法
	Stage      string    `json:"stage"`                  // 阶段
	StatusCode int       `json:"statusCode"`             // 状态码
	RuleID     *string   `json:"ruleId"`                 // 规则ID
	Error      string    `json:"error"`                  // 错误信息
	Timestamp  int64     `gorm:"index" json:"timestamp"` // 时间戳
	CreatedAt  time.Time `json:"createdAt"`              // 创建时间
}

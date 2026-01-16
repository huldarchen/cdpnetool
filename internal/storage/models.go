package storage

import (
	"time"
)

// Setting 用户设置表
type Setting struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 预定义的设置 Key
const (
	SettingKeyDevToolsURL  = "devtools_url"
	SettingKeyTheme        = "theme"
	SettingKeyWindowBounds = "window_bounds"
	SettingKeyLastConfigID = "last_config_id" // 上次使用的配置 ID
)

// ConfigRecord 配置表（存储规则配置）
type ConfigRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`  // 配置描述
	Version     string    `json:"version"`                       // 配置格式版本
	RulesJSON   string    `gorm:"type:text" json:"rulesJson"`    // JSON 序列化的规则数组
	IsActive    bool      `gorm:"default:false" json:"isActive"` // 是否为激活配置
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// InterceptEventRecord 拦截事件历史表
type InterceptEventRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SessionID  string    `gorm:"index" json:"sessionId"`
	TargetID   string    `json:"targetId"`
	Type       string    `gorm:"index" json:"type"` // matched, rewritten, failed, rejected...
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	Stage      string    `json:"stage"`
	StatusCode int       `json:"statusCode"`
	RuleID     *string   `json:"ruleId"`
	Error      string    `json:"error"`
	Timestamp  int64     `gorm:"index" json:"timestamp"`
	CreatedAt  time.Time `json:"createdAt"`
}

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cdpnetool/pkg/rulespec"

	"gorm.io/gorm"
)

// ConfigRepo 配置仓库
type ConfigRepo struct {
	db *DB
}

// NewConfigRepo 创建配置仓库实例
func NewConfigRepo(db *DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}

// Create 创建新配置
func (r *ConfigRepo) Create(name, description, version string, rules []rulespec.Rule) (*ConfigRecord, error) {
	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("序列化规则失败: %w", err)
	}

	record := &ConfigRecord{
		Name:        name,
		Description: description,
		Version:     version,
		RulesJSON:   string(rulesJSON),
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.db.GormDB().Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// Update 更新配置
func (r *ConfigRepo) Update(id uint, name, description, version string, rules []rulespec.Rule) error {
	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("序列化规则失败: %w", err)
	}

	return r.db.GormDB().Model(&ConfigRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"version":     version,
		"rules_json":  string(rulesJSON),
		"updated_at":  time.Now(),
	}).Error
}

// Delete 删除配置
func (r *ConfigRepo) Delete(id uint) error {
	return r.db.GormDB().Delete(&ConfigRecord{}, id).Error
}

// GetByID 根据 ID 获取配置
func (r *ConfigRepo) GetByID(id uint) (*ConfigRecord, error) {
	var record ConfigRecord
	if err := r.db.GormDB().First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetByName 根据名称获取配置
func (r *ConfigRepo) GetByName(name string) (*ConfigRecord, error) {
	var record ConfigRecord
	if err := r.db.GormDB().Where("name = ?", name).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// List 列出所有配置
func (r *ConfigRepo) List() ([]ConfigRecord, error) {
	var records []ConfigRecord
	if err := r.db.GormDB().Order("updated_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// SetActive 设置激活的配置（只能有一个激活）
func (r *ConfigRepo) SetActive(id uint) error {
	return r.db.GormDB().Transaction(func(tx *gorm.DB) error {
		// 先取消所有激活
		if err := tx.Model(&ConfigRecord{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}
		// 激活指定配置
		if err := tx.Model(&ConfigRecord{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetActive 获取当前激活的配置
func (r *ConfigRepo) GetActive() (*ConfigRecord, error) {
	var record ConfigRecord
	if err := r.db.GormDB().Where("is_active = ?", true).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// ParseRules 从记录中解析规则
func (r *ConfigRepo) ParseRules(record *ConfigRecord) ([]rulespec.Rule, error) {
	if record == nil || record.RulesJSON == "" {
		return nil, nil
	}

	var rules []rulespec.Rule
	if err := json.Unmarshal([]byte(record.RulesJSON), &rules); err != nil {
		return nil, fmt.Errorf("解析规则失败: %w", err)
	}
	return rules, nil
}

// ToRulespecConfig 将记录转换为 rulespec.Config
func (r *ConfigRepo) ToRulespecConfig(record *ConfigRecord) (*rulespec.Config, error) {
	rules, err := r.ParseRules(record)
	if err != nil {
		return nil, err
	}

	return &rulespec.Config{
		ID:          fmt.Sprintf("%d", record.ID),
		Name:        record.Name,
		Description: record.Description,
		Version:     record.Version,
		Rules:       rules,
	}, nil
}

// SaveFromRulespecConfig 从 rulespec.Config 保存（更新或创建）
func (r *ConfigRepo) SaveFromRulespecConfig(id uint, name, description string, cfg *rulespec.Config) (*ConfigRecord, error) {
	if id == 0 {
		// 创建新记录
		return r.Create(name, description, cfg.Version, cfg.Rules)
	}
	// 更新现有记录
	if err := r.Update(id, name, description, cfg.Version, cfg.Rules); err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// Rename 重命名配置
func (r *ConfigRepo) Rename(id uint, newName string) error {
	return r.db.GormDB().Model(&ConfigRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":       newName,
		"updated_at": time.Now(),
	}).Error
}

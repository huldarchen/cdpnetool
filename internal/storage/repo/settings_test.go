package repo_test

import (
	"testing"

	"cdpnetool/internal/storage/db"
	"cdpnetool/internal/storage/model"
	"cdpnetool/internal/storage/repo"
)

// setupSettingsTestDB 创建用于 SettingsRepo 测试的内存数据库。
func setupSettingsTestDB(t *testing.T) *repo.SettingsRepo {
	gdb, err := db.New(db.Options{
		FullPath: ":memory:",
		Prefix:   "test_",
	})
	if err != nil {
		t.Fatalf("创建内存数据库失败: %v", err)
	}

	err = db.Migrate(gdb, &model.Setting{})
	if err != nil {
		t.Fatalf("迁移数据库失败: %v", err)
	}

	return repo.NewSettingsRepo(gdb)
}

// TestSettingsRepo_SetAndGet 测试设置的保存与读取。
func TestSettingsRepo_SetAndGet(t *testing.T) {
	r := setupSettingsTestDB(t)

	key := "test_key"
	value := "test_value"

	err := r.Set(key, value)
	if err != nil {
		t.Fatalf("设置失败: %v", err)
	}

	retrieved, err := r.Get(key)
	if err != nil {
		t.Fatalf("获取设置失败: %v", err)
	}

	if retrieved != value {
		t.Errorf("预期值为 %s，实际为 %s", value, retrieved)
	}
}

// TestSettingsRepo_GetWithDefault 测试不存在的键返回默认值。
func TestSettingsRepo_GetWithDefault(t *testing.T) {
	r := setupSettingsTestDB(t)

	defaultVal := "default_value"
	retrieved := r.GetWithDefault("non_existent_key", defaultVal)

	if retrieved != defaultVal {
		t.Errorf("预期返回默认值 %s，实际返回 %s", defaultVal, retrieved)
	}
}

// TestSettingsRepo_SetMultiple 测试批量设置功能及事务一致性。
func TestSettingsRepo_SetMultiple(t *testing.T) {
	r := setupSettingsTestDB(t)

	kvs := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	err := r.SetMultiple(kvs)
	if err != nil {
		t.Fatalf("批量设置失败: %v", err)
	}

	// 验证所有键值对是否正确保存
	for key, expectedVal := range kvs {
		actualVal, err := r.Get(key)
		if err != nil {
			t.Errorf("获取键 %s 失败: %v", key, err)
		}
		if actualVal != expectedVal {
			t.Errorf("键 %s 预期值 %s，实际值 %s", key, expectedVal, actualVal)
		}
	}
}

// TestSettingsRepo_DeleteByKey 测试按键删除功能。
func TestSettingsRepo_DeleteByKey(t *testing.T) {
	r := setupSettingsTestDB(t)

	key := "to_delete"
	r.Set(key, "some_value")

	err := r.DeleteByKey(key)
	if err != nil {
		t.Fatalf("删除失败: %v", err)
	}

	_, err = r.Get(key)
	if err == nil {
		t.Error("预期键已删除，但仍然能获取到值")
	}
}

// TestSettingsRepo_PresetKeys 测试预设置的键是否按预期工作。
func TestSettingsRepo_PresetKeys(t *testing.T) {
	r := setupSettingsTestDB(t)

	// 测试 DevToolsURL
	expectedURL := "http://localhost:9999"
	r.SetDevToolsURL(expectedURL)
	actualURL := r.GetDevToolsURL()
	if actualURL != expectedURL {
		t.Errorf("DevToolsURL 预期 %s，实际 %s", expectedURL, actualURL)
	}

	// 测试 Theme
	expectedTheme := "dark"
	r.SetTheme(expectedTheme)
	actualTheme := r.GetTheme()
	if actualTheme != expectedTheme {
		t.Errorf("Theme 预期 %s，实际 %s", expectedTheme, actualTheme)
	}

	// 测试默认值
	defaultURL := r.GetDevToolsURL()
	r.DeleteByKey(model.SettingKeyDevToolsURL)
	resetURL := r.GetDevToolsURL()
	if resetURL == defaultURL {
		// 应该返回默认值 "http://localhost:9222"
		if resetURL != "http://localhost:9222" {
			t.Errorf("DevToolsURL 默认值不符合预期，实际为 %s", resetURL)
		}
	}
}

package repo_test

import (
	"context"
	"testing"

	"cdpnetool/internal/storage/db"
	"cdpnetool/internal/storage/model"
	"cdpnetool/internal/storage/repo"
	"cdpnetool/pkg/rulespec"
)

// setupTestDB 创建一个用于测试的内存数据库并完成迁移。
func setupTestDB(t *testing.T) *repo.ConfigRepo {
	// 使用内存数据库，速度快且隔离性好
	gdb, err := db.New(db.Options{
		Name:   ":memory:",
		Prefix: "test_",
	})
	if err != nil {
		t.Fatalf("创建内存数据库失败: %v", err)
	}

	// 执行必要迁移
	err = db.Migrate(gdb, &model.ConfigRecord{})
	if err != nil {
		t.Fatalf("迁移数据库失败: %v", err)
	}

	return repo.NewConfigRepo(gdb)
}

// TestConfigRepo_Create 测试配置的创建功能。
// 验证配置 ID 校验、规则 ID 校验以及数据持久化。
func TestConfigRepo_Create(t *testing.T) {
	r := setupTestDB(t)

	t.Run("Create Valid Config", func(t *testing.T) {
		cfg := rulespec.NewConfig("测试配置")
		record, err := r.Create(context.Background(), cfg)
		if err != nil {
			t.Fatalf("创建配置失败: %v", err)
		}

		if record.Name != "测试配置" {
			t.Errorf("预期名称为 '测试配置'，实际为 %s", record.Name)
		}
		if record.ConfigID != cfg.ID {
			t.Errorf("预期 ConfigID 为 %s，实际为 %s", cfg.ID, record.ConfigID)
		}
	})

	t.Run("Create Duplicate Rule IDs", func(t *testing.T) {
		cfg := rulespec.NewConfig("重复规则测试")
		cfg.Rules = []rulespec.Rule{
			{ID: "rule1", Name: "R1"},
			{ID: "rule1", Name: "R2"}, // 重复 ID
		}
		_, err := r.Create(context.Background(), cfg)
		if err == nil {
			t.Error("预期因规则 ID 重复而报错，但实际未报错")
		}
	})
}

// TestConfigRepo_SetActive 测试激活配置的切换逻辑。
// 验证系统是否能确保同时只有一个配置处于激活状态。
func TestConfigRepo_SetActive(t *testing.T) {
	r := setupTestDB(t)

	// 创建两个配置
	c1, _ := r.Create(context.Background(), rulespec.NewConfig("C1"))
	c2, _ := r.Create(context.Background(), rulespec.NewConfig("C2"))

	// 1. 激活第一个
	err := r.SetActive(context.Background(), c1.ID)
	if err != nil {
		t.Fatalf("激活 C1 失败: %v", err)
	}

	active, _ := r.GetActive(context.Background())
	if active.ID != c1.ID {
		t.Errorf("预期激活 ID 为 %d，实际为 %d", c1.ID, active.ID)
	}

	// 2. 激活第二个，验证第一个被取消激活
	err = r.SetActive(context.Background(), c2.ID)
	if err != nil {
		t.Fatalf("激活 C2 失败: %v", err)
	}

	active, _ = r.GetActive(context.Background())
	if active.ID != c2.ID {
		t.Errorf("预期激活 ID 为 %d，实际为 %d", c2.ID, active.ID)
	}

	// 3. 验证 C1 已不活跃
	var c1Record model.ConfigRecord
	r.Db.First(&c1Record, c1.ID)
	if c1Record.IsActive {
		t.Error("C1 应该处于非激活状态")
	}
}

// TestConfigRepo_Upsert 测试导入配置时的更新或创建逻辑。
func TestConfigRepo_Upsert(t *testing.T) {
	r := setupTestDB(t)
	cfg := rulespec.NewConfig("初始配置")

	// 1. 第一次 Upsert 应为创建
	record, err := r.Upsert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("首次 Upsert 失败: %v", err)
	}
	initialID := record.ID

	// 2. 修改名称后再次 Upsert (ID 相同)
	cfg.Name = "已更新名称"
	updated, err := r.Upsert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("更新 Upsert 失败: %v", err)
	}

	if updated.ID != initialID {
		t.Errorf("预期更新后的数据库 ID 保持不变，实际 %d -> %d", initialID, updated.ID)
	}
	if updated.Name != "已更新名称" {
		t.Errorf("预期名称已更新为 '已更新名称'，实际为 %s", updated.Name)
	}
}

// TestConfigRepo_Rename 测试重命名功能。
// 验证数据库记录名称和内部 JSON 字符串中的名称同步更新。
func TestConfigRepo_Rename(t *testing.T) {
	r := setupTestDB(t)
	cfg := rulespec.NewConfig("原名")
	record, _ := r.Create(context.Background(), cfg)

	newName := "新名"
	err := r.Rename(context.Background(), record.ID, newName)
	if err != nil {
		t.Fatalf("重命名失败: %v", err)
	}

	// 重新获取并验证
	updated, _ := r.FindOne(context.Background(), record.ID)
	if updated.Name != newName {
		t.Errorf("数据库记录名称未更新，预期 %s，实际 %s", newName, updated.Name)
	}

	// 验证 JSON 内部是否也更新了
	parsed, _ := r.ToRulespecConfig(updated)
	if parsed.Name != newName {
		t.Errorf("配置 JSON 内部名称未更新，预期 %s，实际 %s", newName, parsed.Name)
	}
}

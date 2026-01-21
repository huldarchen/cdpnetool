package db_test

import (
	"cdpnetool/internal/storage/db"
	"path/filepath"
	"strings"
	"testing"
)

// TestModel 定义一个用于测试数据库迁移和基础操作的简单模型。
type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

// TestGetDefaultPath 测试数据库默认存储路径的生成逻辑。
// 验证路径是否包含应用名称，并以指定的数据库文件名结尾。
func TestGetDefaultPath(t *testing.T) {
	t.Run("Generate Path", func(t *testing.T) {
		dbName := "test_db.db"
		path, err := db.GetDefaultPath(dbName)
		if err != nil {
			t.Fatalf("获取默认路径失败: %v", err)
		}

		if !strings.HasSuffix(path, dbName) {
			t.Errorf("路径 %s 不是以 %s 结尾", path, dbName)
		}

		if !strings.Contains(path, "cdpnetool") {
			t.Errorf("路径 %s 不包含应用名称 'cdpnetool'", path)
		}
	})
}

// TestDatabaseInitialization 测试数据库的初始化、连接以及自动迁移功能。
// 验证表前缀配置、SingularTable 策略以及基本的数据读写。
func TestDatabaseInitialization(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "unit_test.db")

	opts := db.Options{
		FullPath: dbPath,
		Prefix:   "test_",
	}

	// 1. 测试数据库连接创建
	gdb, err := db.New(opts)
	if err != nil {
		t.Fatalf("初始化数据库连接失败: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("获取底层的 sql.DB 失败: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Errorf("数据库 Ping 失败: %v", err)
	}

	// 2. 测试自动迁移功能
	err = db.Migrate(gdb, &TestModel{})
	if err != nil {
		t.Fatalf("执行数据库迁移失败: %v", err)
	}

	// 3. 测试数据写入和表前缀验证
	testEntry := TestModel{Name: "clean_code"}
	if err := gdb.Create(&testEntry).Error; err != nil {
		t.Errorf("向迁移后的表中写入数据失败: %v", err)
	}

	var count int64
	// 验证记录数，同时也验证了模型与表的匹配
	if err := gdb.Model(&TestModel{}).Count(&count).Error; err != nil {
		t.Errorf("查询记录数失败: %v", err)
	}

	if count != 1 {
		t.Errorf("预期记录数为 1，实际为 %d", count)
	}

	// 验证物理表名是否符合前缀规则
	var tableName string
	// SQLite 特有的查询表名方式
	row := sqlDB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test_test_model'")
	if err := row.Scan(&tableName); err != nil {
		t.Errorf("未找到预期的带前缀表名 'test_test_model': %v", err)
	}
}

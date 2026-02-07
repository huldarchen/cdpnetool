package db

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Options 数据库配置选项
type Options struct {
	// Name 数据库文件名
	Name string
	// Dir 数据库目录路径（如果指定，则优先使用此目录；否则使用默认目录）
	Dir string
	// Prefix 表前缀
	Prefix string
	// Logger GORM 日志实现
	Logger logger.Interface
}

// New 创建并初始化数据库连接。
// 它会根据 Options 中提供的路径信息打开 SQLite 数据库，并配置命名策略和连接池。
func New(opts Options) (*gorm.DB, error) {
	var dbPath string

	// 处理 SQLite 内存数据库特殊情况
	if opts.Name == ":memory:" {
		dbPath = ":memory:"
	} else {
		// 获取数据库目录
		dbDir := opts.Dir
		if dbDir == "" {
			var err error
			dbDir, err = GetDefaultDir()
			if err != nil {
				return nil, err
			}
		}

		// 确保数据库目录存在
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, err
		}

		// 拼接完整的数据库文件路径
		dbPath = filepath.Join(dbDir, opts.Name)
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: opts.Logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   opts.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池（对于 SQLite 主要是控制并发）
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}

	return db, nil
}

// Migrate 执行数据库自动迁移
func Migrate(db *gorm.DB, models ...any) error {
	return db.AutoMigrate(models...)
}

// GetDefaultDir 获取平台相关的默认数据库目录（不包含文件名）
func GetDefaultDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// %APPDATA%/cdpnetool/
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	case "darwin":
		// ~/Library/Application Support/cdpnetool/
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, "Library", "Application Support")
	default:
		// Linux: ~/.local/share/cdpnetool/
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			baseDir = filepath.Join(home, ".local", "share")
		}
	}

	return filepath.Join(baseDir, "cdpnetool"), nil
}

package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	dbType string
	config gdb.ConfigNode
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

// Initialize 初始化数据库连接
func (dm *DatabaseManager) Initialize(ctx context.Context) error {
	cfg := g.Cfg()

	// 获取数据库类型
	dbType := cfg.MustGet(ctx, "database.type", "mysql").String()
	dm.dbType = dbType

	g.Log().Infof(ctx, "正在初始化数据库类型: %s", dbType)

	// 根据数据库类型获取配置
	var configKey string
	switch dbType {
	case "mysql":
		configKey = "database.mysql"
	case "sqlite":
		configKey = "database.sqlite"
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbType)
	}

	// 读取对应的数据库配置
	var config gdb.ConfigNode
	if err := cfg.MustGet(ctx, configKey).Struct(&config); err != nil {
		return fmt.Errorf("读取数据库配置失败: %v", err)
	}

	// 确保配置中的 Type 字段被正确设置
	config.Type = dbType
	dm.config = config

	// 如果是 SQLite，确保数据目录存在
	if dbType == "sqlite" {
		if err := dm.ensureSQLiteDir(); err != nil {
			return fmt.Errorf("创建 SQLite 数据目录失败: %v", err)
		}
	}

	// 设置数据库配置
	err := gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{config},
	})
	if err != nil {
		return err
	}

	// 测试数据库连接
	if err := dm.testConnection(ctx); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	g.Log().Infof(ctx, "数据库初始化成功: %s", dbType)
	return nil
}

// ensureSQLiteDir 确保 SQLite 数据目录存在
func (dm *DatabaseManager) ensureSQLiteDir() error {
	// 从连接字符串中提取数据库文件路径
	link := dm.config.Link
	if link == "" {
		return fmt.Errorf("SQLite 连接字符串为空")
	}

	g.Log().Infof(context.Background(), "正在解析 SQLite 连接字符串: %s", link)

	var dbPath string

	// 解析不同格式的 SQLite 连接字符串
	if len(link) > 7 && link[:7] == "sqlite:" {
		remaining := link[7:] // 去除 "sqlite:" 前缀

		// 处理不同的连接字符串格式
		if remaining[:6] == "@file(" && remaining[len(remaining)-1] == ')' {
			// 格式: sqlite:@file(./data/omniscient.db)
			dbPath = remaining[6 : len(remaining)-1]
		} else if remaining[0] == ':' && len(remaining) > 6 && remaining[1:7] == "@file(" && remaining[len(remaining)-1] == ')' {
			// 格式: sqlite::@file(./data/omniscient.db)
			dbPath = remaining[7 : len(remaining)-1]
		} else {
			// 简单格式: sqlite:./data/omniscient.db
			dbPath = remaining
		}
	} else {
		// 如果不是以 sqlite: 开头，直接使用原字符串
		dbPath = link
	}

	g.Log().Infof(context.Background(), "提取的数据库文件路径: %s", dbPath)

	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(dbPath) {
		workDir, _ := os.Getwd()
		dbPath = filepath.Join(workDir, dbPath)
	}

	// 确保目录存在
	dbDir := filepath.Dir(dbPath)
	g.Log().Infof(context.Background(), "数据库目录路径: %s", dbDir)

	if !gfile.Exists(dbDir) {
		if err := gfile.Mkdir(dbDir); err != nil {
			return fmt.Errorf("创建数据库目录失败 %s: %v", dbDir, err)
		}
		g.Log().Infof(context.Background(), "创建 SQLite 数据目录: %s", dbDir)
	} else {
		g.Log().Infof(context.Background(), "SQLite 数据目录已存在: %s", dbDir)
	}

	return nil
}

// testConnection 测试数据库连接
func (dm *DatabaseManager) testConnection(ctx context.Context) error {
	db := g.DB()

	var testSQL string
	switch dm.dbType {
	case "mysql":
		testSQL = "SELECT 1"
	case "sqlite":
		testSQL = "SELECT 1"
	default:
		testSQL = "SELECT 1"
	}

	_, err := db.GetValue(ctx, testSQL)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	g.Log().Info(ctx, "数据库连接测试成功")
	return nil
}

// GetDatabaseType 获取当前数据库类型
func (dm *DatabaseManager) GetDatabaseType() string {
	return dm.dbType
}

// tableExists 检查表是否存在
func (dm *DatabaseManager) tableExists(ctx context.Context, tableName string) (bool, error) {
	db := g.DB()

	var checkSQL string
	switch dm.dbType {
	case "mysql":
		checkSQL = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	case "sqlite":
		checkSQL = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?"
	default:
		return false, fmt.Errorf("不支持的数据库类型: %s", dm.dbType)
	}

	count, err := db.GetValue(ctx, checkSQL, tableName)
	if err != nil {
		return false, err
	}

	return count.Int() > 0, nil
}

// CreateTables 创建数据表（避免重复创建）
func (dm *DatabaseManager) CreateTables(ctx context.Context) error {
	db := g.DB()

	// 检查 jpid 表是否已存在
	exists, err := dm.tableExists(ctx, "jpid")
	if err != nil {
		g.Log().Warningf(ctx, "检查表是否存在时出错: %v", err)
	}

	if exists {
		g.Log().Info(ctx, "数据表 jpid 已存在，跳过创建")
		return nil
	}

	g.Log().Info(ctx, "正在创建数据表 jpid...")

	// 根据数据库类型选择合适的 SQL
	var createJpidTableSQL string
	switch dm.dbType {
	case "mysql":
		createJpidTableSQL = `
			CREATE TABLE IF NOT EXISTS jpid (
				id INT NOT NULL AUTO_INCREMENT,
				name VARCHAR(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'java项目名',
				ports VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '运行端口,多个逗号隔开',
				pid INT NOT NULL COMMENT 'pid',
				catalog VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '运行目录',
				run LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '原生启动命令',
				script LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'sh脚本启动命令',
				worker VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '服务器',
				status INT DEFAULT '0' COMMENT '状态[1:启动，0:停止]',
				description VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '项目描述',
				way INT DEFAULT '2' COMMENT '启动方式[1:docker, 2:jdk]',
				autostart INT DEFAULT '0' COMMENT '自启[0:没有自启, 1:自启]',
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='java项目详情';
			`
	case "sqlite":
		createJpidTableSQL = `
			CREATE TABLE IF NOT EXISTS jpid (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL, -- java项目名
				ports TEXT NOT NULL, -- 运行端口,多个逗号隔开
				pid INTEGER NOT NULL, -- pid
				catalog TEXT DEFAULT NULL, -- 运行目录
				run TEXT, -- 原生启动命令
				script TEXT, -- sh脚本启动命令
				worker TEXT NOT NULL, -- 服务器
				status INTEGER DEFAULT 0, -- 状态[1:启动，0:停止]
				description TEXT DEFAULT NULL, -- 项目描述
				way INTEGER DEFAULT 2, -- 启动方式[1:docker, 2:jdk]
				autostart INTEGER DEFAULT 0 -- 自启[0:没有自启, 1:自启]
			);
	   `
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dm.dbType)
	}

	// 执行建表语句
	_, err = db.Exec(ctx, createJpidTableSQL)
	if err != nil {
		return fmt.Errorf("创建 jpid 表失败: %v", err)
	}

	g.Log().Info(ctx, "数据表 jpid 创建成功")
	return nil
}

// GetDatabaseInfo 获取数据库信息
func (dm *DatabaseManager) GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error) {
	db := g.DB()
	info := make(map[string]interface{})

	info["type"] = dm.dbType
	info["link"] = dm.config.Link

	switch dm.dbType {
	case "mysql":
		// 获取 MySQL 版本
		version, err := db.GetValue(ctx, "SELECT VERSION()")
		if err == nil {
			info["version"] = version
		}

		// 获取数据库名
		dbName, err := db.GetValue(ctx, "SELECT DATABASE()")
		if err == nil {
			info["database"] = dbName
		}

		// 获取表信息
		tables, err := db.GetAll(ctx, "SHOW TABLES")
		if err == nil {
			tableNames := make([]string, 0)
			for _, table := range tables {
				for _, v := range table {
					tableNames = append(tableNames, v.String())
				}
			}
			info["tables"] = tableNames
		}

	case "sqlite":
		// 获取 SQLite 版本
		version, err := db.GetValue(ctx, "SELECT sqlite_version()")
		if err == nil {
			info["version"] = version
		}

		// 解析数据库文件路径
		link := dm.config.Link
		var dbPath string
		if len(link) > 7 && link[:7] == "sqlite:" {
			remaining := link[7:]
			if remaining[:6] == "@file(" && remaining[len(remaining)-1] == ')' {
				dbPath = remaining[6 : len(remaining)-1]
			} else if remaining[0] == ':' && len(remaining) > 6 && remaining[1:7] == "@file(" && remaining[len(remaining)-1] == ')' {
				dbPath = remaining[7 : len(remaining)-1]
			} else {
				dbPath = remaining
			}
		} else {
			dbPath = link
		}

		if dbPath != "" {
			info["database_file"] = dbPath
		}

		// 获取表信息
		tables, err := db.GetAll(ctx, "SELECT name FROM sqlite_master WHERE type='table'")
		if err == nil {
			tableNames := make([]string, 0)
			for _, table := range tables {
				tableNames = append(tableNames, table["name"].String())
			}
			info["tables"] = tableNames
		}
	}

	return info, nil
}

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

	// 解析 SQLite 路径 (格式: sqlite:./data/omniscient.db)
	if len(link) > 7 && link[:7] == "sqlite:" {
		dbPath := link[7:] // 去除 "sqlite:" 前缀

		// 如果是相对路径，转换为绝对路径
		if !filepath.IsAbs(dbPath) {
			workDir, _ := os.Getwd()
			dbPath = filepath.Join(workDir, dbPath)
		}

		// 确保目录存在
		dbDir := filepath.Dir(dbPath)
		if !gfile.Exists(dbDir) {
			if err := gfile.Mkdir(dbDir); err != nil {
				return fmt.Errorf("创建数据库目录失败 %s: %v", dbDir, err)
			}
			g.Log().Infof(context.Background(), "创建 SQLite 数据目录: %s", dbDir)
		}
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

// CreateTables 创建数据表（示例）
func (dm *DatabaseManager) CreateTables(ctx context.Context) error {
	db := g.DB()

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
	_, err := db.Exec(ctx, createJpidTableSQL)
	if err != nil {
		return fmt.Errorf("创建用户表失败: %v", err)
	}

	g.Log().Info(ctx, "数据表创建成功")
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

	case "sqlite":
		// 获取 SQLite 版本
		version, err := db.GetValue(ctx, "SELECT sqlite_version()")
		if err == nil {
			info["version"] = version
		}

		// 数据库文件路径
		if len(dm.config.Link) > 7 {
			info["database_file"] = dm.config.Link[7:]
		}
	}

	return info, nil
}

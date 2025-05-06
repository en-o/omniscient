// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// JpidDao is the data access object for the table jpid.
type JpidDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  JpidColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// JpidColumns defines and stores column names for the table jpid.
type JpidColumns struct {
	Id          string //
	Name        string // java项目名
	Port        string // 运行端口
	Pid         string // pid
	Catalog     string // 运行目录
	Run         string // 运行脚本（sh命令
	Description string // 项目描述
}

// jpidColumns holds the columns for the table jpid.
var jpidColumns = JpidColumns{
	Id:          "id",
	Name:        "name",
	Port:        "port",
	Pid:         "pid",
	Catalog:     "catalog",
	Run:         "run",
	Description: "description",
}

// NewJpidDao creates and returns a new DAO object for table data access.
func NewJpidDao(handlers ...gdb.ModelHandler) *JpidDao {
	return &JpidDao{
		group:    "default",
		table:    "jpid",
		columns:  jpidColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *JpidDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *JpidDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *JpidDao) Columns() JpidColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *JpidDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *JpidDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *JpidDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

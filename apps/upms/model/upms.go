package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/yanshicheng/ikube-gin-xjob/common/model"
	"gorm.io/gorm"
)

// 角色表，角色菜单关联表，角色账户关联表， 权限表，
func init() {
	model.Register(&Menu{}, &Role{}, &RoleMenu{}, &Upms{})
}

const MenuLevel = 3

type Menu struct {
	model.Model
	Path             string  `json:"path" form:"path" binding:"required,max=32" gorm:"type:varchar(32);not null;comment:路由路径"`
	Name             string  `json:"name" form:"name" binding:"required,max=32" gorm:"type:varchar(32);not null;unique;comment:唯一标识名称" `
	Component        string  `json:"component" form:"component" binding:"required,max=255" gorm:"type:varchar(255);not null;comment:组件路径" `
	Redirect         string  `json:"redirect" form:"redirect" binding:"max=255" gorm:"type:varchar(255);comment:重定向路径" `
	Title            string  `json:"title" form:"title" binding:"max=26" gorm:"type:varchar(26);not null;comment:菜单标题" `
	Icon             string  `json:"icon"  form:"icon" binding:"max=32" gorm:"type:varchar(32);comment:菜单图标" `
	Expanded         bool    `json:"expanded"  form:"expanded" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否默认展开" `
	OrderNo          int     `json:"orderNo" form:"orderNo" binding:"required,number" gorm:"type:tinyint;not null;comment:菜单顺序编号" `
	Hidden           bool    `json:"hidden"  form:"hidden" binding:"required" gorm:"default:false;comment:是否隐藏菜单"`
	HiddenBreadcrumb bool    `json:"hiddenBreadcrumb" form:"hiddenBreadcrumb" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否隐藏面包屑"`
	Single           bool    `json:"single" form:"single" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否单级菜单显示"`
	FrameSrc         string  `json:"frameSrc" form:"frameSrc" binding:"max=255" gorm:"type:varchar(255);comment:内嵌iframe的地址"`
	FrameBlank       bool    `json:"frameBlank" form:"frameBlank" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:内嵌iframe是否新窗口打开" `
	KeepAlive        bool    `json:"keepAlive" form:"keepAlive" binding:"boolean" gorm:"type:tinyint(1);default:true;comment:开启keep-alive"`
	ParentId         uint    `json:"parentId" form:"parentId"  binding:"required,number" gorm:"type:int;not null;comment:父级"` // 关联父级路由
	Level            int     `json:"level" form:"level" gorm:"type:int;not null;comment:层级"`
	Children         []*Menu `gorm:"-" json:"children"` // 子路由，不存储在数据库中，只用于加载和显示

}

func (*Menu) TableName() string {
	return "ikubexjob_upms_menu"
}

// BeforeCreate 机构表 创建钩子函数
func (o *Menu) BeforeCreate(tx *gorm.DB) error {
	// 检查是否有父节点，如果没有父节点，则为根节点
	if o.ParentId != 0 {
		// 如果 ParentID 不为0，说明此节点有父节点

		var parent Menu
		// 查询父节点的详细信息
		// 这里使用 tx.First 来查询具有指定 ID 的父节点
		// o.ParentID 是父节点的 ID，将结果存储在 parent 变量中
		if err := tx.First(&parent, o.ParentId).Error; err != nil {
			// 如果查询过程中出现错误，例如数据库连接错误或找不到指定的父节点
			return err // 返回错误，中断创建操作
		}

		// 如果父节点查询成功，设置当前节点的层级为父节点层级 + 1
		o.Level = parent.Level + 1

		// 检查层级是否超过5
		if o.Level > MenuLevel {
			// 如果层级超过5，返回错误
			return fmt.Errorf("cannot add beyond level %d", MenuLevel)
		}
	} else {
		// 如果 ParentID 为0，说明此节点没有父节点，即它是一个根节点
		o.Level = 1 // 设置根节点的层级为1
	}
	// 如果所有检查都通过，没有错误，则返回 nil，允许创建操作继续进行
	return nil
}

type Role struct {
	model.Model
	Name string `json:"name" form:"name" binding:"required,alphanum,max=32" gorm:"type:varchar(32);not null;unique;comment:角色"`
}

func (r *Role) TableName() string {
	return "ikubexjob_upms_role"
}

type RoleMenu struct {
	model.Model
	RoleId uint `json:"roleId" form:"roleId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_menu;comment:角色" `
	MenuId uint `json:"menuId" form:"menuId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_menu;comment:菜单" `
}

func (r *RoleMenu) TableName() string {
	return "ikubexjob_upms_role_menu"
}

type Upms struct {
	model.Model
	Name     string     `json:"name" form:"name" binding:"required,max=32" gorm:"type:varchar(32);not null;unique;comment:权限名称"`
	RoleId   uint       `json:"roleId" form:"roleId" binding:"required,number" gorm:"type:int;not null;comment:角色"`
	Resource string     `json:"resource" form:"resource" binding:"required,max=255" gorm:"type:varchar(255);not null;comment:资源"`
	Type     ActionType `json:"type" form:"type"  binding:"required,oneof=0 1" gorm:"type:tinyint;not null;comment:操作类型"`
}

func (u *Upms) TableName() string {
	return "ikubexjob_upms_upms"
}

// ActionType  定义 ActionType 类型
type ActionType uint

const (
	ReadAction  ActionType = 0
	WriteAction ActionType = 1
)

// ActionTypeToString 将 ActionType 转换为字符串

func (a *ActionType) String() string {
	switch *a {
	case ReadAction:
		return "read"
	case WriteAction:
		return "write"
	default:
		return "unknown"
	}
}

// Scan 实现  接口
func (a *ActionType) Scan(value interface{}) error {
	v, ok := value.(uint)
	if !ok {
		return fmt.Errorf("invalid value for ActionType: %v", value)
	}
	*a = ActionType(v)
	return nil
}

// Value 实现 Valuer 接口
func (a *ActionType) Value() (driver.Value, error) {
	return uint(*a), nil
}

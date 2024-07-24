package types

import (
	"github.com/yanshicheng/ikube-gin-xjob/common/model"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"time"
)

type AccountLoginReq struct {
	Account  string `json:"account" form:"account" binding:"required,max=32"`
	Password string `json:"password" form:"password" binding:"required,max=128"`
}

type AccountChangePasswordReq struct {
	Account       string `json:"account" form:"account" binding:"required,max=32"`
	Password      string `json:"password" form:"password" binding:"required,max=128"`
	NewPassword   string `json:"newPassword" form:"newPassword" binding:"required,max=128"`
	ReNewPassword string `json:"reNewPassword" form:"reNewPassword" binding:"required,max=128,eqfield=NewPassword"`
}

type AccountRestPasswordReq struct {
	Account string `json:"account" form:"account" binding:"required,max=32"`
}

type AccountIconResp struct {
	IconPath string `json:"iconPath"`
}

type AccountQueryReq struct {
	types.Pagination
	UserName       string `json:"userName" form:"userName" uri:"userName"`
	Email          string `json:"email" form:"email" uri:"email"`
	WorkNumber     string `json:"workNumber" form:"workNumber" uri:"workNumber"`
	Mobile         string `json:"mobile" form:"mobile" uri:"mobile"`
	Account        string `json:"account" form:"account" uri:"account"`
	IsFrozen       bool   `json:"isFrozen" form:"isFrozen"`
	IsDisabled     bool   `json:"isDisabled" form:"isDisabled"`
	IsLeave        bool   `json:"isLeave" form:"isLeave" `
	PositionId     *uint  `json:"PositionId" form:"PositionId" ` // 对应职位表
	OrganizationId *uint  `json:"organizationId" form:"organizationId"`
}

type AccountCreateReq struct {
	UserName       string         `json:"userName" form:"userName" binding:"required,max=32"`
	Account        string         `json:"account" form:"account"  binding:"required,max=32"`
	Mobile         string         `json:"mobile" from:"mobile" binding:"required,max=11"`
	Email          string         `json:"email" form:"email" binding:"required,max=36,email"`
	WorkNumber     string         `json:"workNumber" form:"workNumber" binding:"required,max=24"`
	HireDate       model.DateTime `json:"hireDate" form:"hireDate" binding:"required"  time_format:"2006-01-02"`
	PositionId     uint           `json:"PositionId" form:"PositionId" binding:"required,number" gorm:"type:int;not null;comment:职位Id"` // 对应职位表
	OrganizationId uint           `json:"organizationId" form:"organizationId" binding:"required,number" gorm:"type:int;not null;comment:组织ID"`
}

type AccountListResp struct {
	Id             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt" gorm:"type:datetime;autoCreateTime;comment:创建时间"` // 创建时间
	UpdatedAt      time.Time `json:"updatedAt" gorm:"type:datetime;autoUpdateTime;comment:更新时间"`
	UserName       string    `json:"userName" form:"userName" binding:"required,max=32" gorm:"type:varchar(32);not null;comment:姓名"`
	Account        string    `json:"account" form:"Account" binding:"required,max=32" gorm:"type:varchar(32);unique_index;not null;comment:账号"`
	Password       string    `json:"password" form:"password" binding:"max=24" gorm:"type:varchar(256);not null;comment:密码"`
	Icon           string    `json:"icon" form:"icon" gorm:"type:varchar(256);not null;comment:头像"`
	Mobile         string    `json:"mobile" form:"mobile" binding:"required,max=11" gorm:"type:char(11);unique_index;not null;comment:手机号"`
	Email          string    `json:"email" form:"email" binding:"required,max=36,email" gorm:"type:varchar(36);unique_index;not null;comment:邮箱"`
	WorkNumber     string    `json:"workNumber" form:"workNumber" binding:"required,max=24" gorm:"type:varchar(24);unique_index;not null;comment:工号"`
	HireDate       time.Time `json:"hireDate" form:"hireDate" binding:"required" gorm:"type:date;not null;comment:入职时间"`
	IsFrozen       bool      `json:"isFrozen" form:"isFrozen" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否冻结"`
	IsDisabled     bool      `json:"isDisabled" form:"isDisabled" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否禁用"`
	IsLeave        bool      `json:"isLeave" form:"isLeave" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否离职"`
	PositionId     int       `json:"PositionId" form:"PositionId" binding:"required,number" gorm:"type:int;not null;comment:职位Id"` // 对应职位表
	OrganizationId uint      `json:"organizationId" form:"organizationId" binding:"required,number" gorm:"type:int;not null;comment:组织ID"`
}

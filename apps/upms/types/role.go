package types

import "github.com/yanshicheng/ikube-gin-xjob/common/types"

type RoleSearchReq struct {
	Name string `json:"name" form:"name" uri:"name"`
	types.Pagination
}

type RoleUpdateRequest struct {
	Name string `json:"name" form:"name" binding:"required,max=32"`
}

type RoleAccountBindRequest struct {
	AccountId []uint `json:"accountId" form:"accountId" binding:"required"`
	RoleId    uint   `json:"roleId" form:"roleId" binding:"required,number"`
}

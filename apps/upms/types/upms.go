package types

import "github.com/yanshicheng/ikube-gin-xjob/common/types"

type UpmsSearchReq struct {
	Name   string `json:"name" form:"name" uri:"name"`
	RoleId uint   `json:"roleId" form:"roleId" uri:"roleId"`
	types.Pagination
}

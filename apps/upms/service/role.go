package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
)

type RoleService interface {
	//Get(*gin.Context, types.SearchId) (*model.Role, error)

	List(*gin.Context, types2.RoleSearchReq) (*types.QueryResponse, error)
	Create(*gin.Context, *model.Role) error
	Put(*gin.Context, types.SearchId, *types2.RoleUpdateRequest) (*model.Role, error)
	Delete(*gin.Context, types.SearchId) error
}

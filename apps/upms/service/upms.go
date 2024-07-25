package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
)

type UpmsService interface {
	//Get(*gin.Context, types.SearchId) (*model.Role, error)
	List(*gin.Context, types2.UpmsSearchReq) (*types.QueryResponse, error)
	Create(*gin.Context, *model.Upms) error
	Put(*gin.Context, types.SearchId, *model.Upms) (*model.Upms, error)
	Delete(*gin.Context, types.SearchId) error
}

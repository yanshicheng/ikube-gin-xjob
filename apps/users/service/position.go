package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
)

type PositionService interface {
	List(*gin.Context, types2.PositionListSearchReq) ([]*model.Position, error)
	Create(*gin.Context, *model.Position) error
	Put(*gin.Context, types.SearchId, *model.Position) (*model.Position, error)
	Delete(*gin.Context, types.SearchId) error
}

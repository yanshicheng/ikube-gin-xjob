package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
)

type MenuService interface {
	//Get(*gin.Context, types.SearchId) (*model.Role, error)
	List(*gin.Context, types2.MenuSearchReq) ([]*model.Menu, error)
	Create(*gin.Context, *model.Menu) error
	Put(*gin.Context, types.SearchId, *model.Menu) (*model.Menu, error)
	Delete(*gin.Context, types.SearchId) error
}

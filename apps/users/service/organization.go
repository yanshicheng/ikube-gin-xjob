package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	otypes "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
)

type OrganizationService interface {
	Get(*gin.Context, otypes.OrganizationGetSearchReq) ([]*model.Organization, error)
	List(*gin.Context, otypes.OrganizationGetSearchReq) ([]*model.Organization, error)
	Create(*gin.Context, *model.Organization) error
	Put(*gin.Context, types.SearchId, *model.Organization) (*model.Organization, error)
	Delete(*gin.Context, types.SearchId) error
}

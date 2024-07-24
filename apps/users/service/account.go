package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/errorx"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/utils"
)

type AccountService interface {
	Create(*gin.Context, *types2.AccountCreateReq) (*model.Account, error)
	Delete(*gin.Context, types.SearchId) error
	Get(*gin.Context, types.SearchId) (*model.Account, error)
	List(*gin.Context, types2.AccountQueryReq) (*types.QueryResponse, error) // TODO 分页
	Put(*gin.Context, types.SearchId, *types2.AccountCreateReq) (*model.Account, error)
	RestPassword(*gin.Context, *types2.AccountRestPasswordReq) error
	ChangePassword(*gin.Context, *types2.AccountChangePasswordReq) error
	Login(*gin.Context, *types2.AccountLoginReq) (*utils.JWTResponse, errorx.ErrorCode, error)
	Logout(*gin.Context) error
	ChangeIcon(*gin.Context) (types2.AccountIconResp, error)
}

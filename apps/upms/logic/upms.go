package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/upms"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/service"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ service.UpmsService = (*UpmsHandler)(nil)

type UpmsHandler struct {
	l  *zap.Logger
	db *gorm.DB
}

func (l *UpmsHandler) List(ctx *gin.Context, req types2.UpmsSearchReq) (*types.QueryResponse, error) {

	return nil, nil
}

func (l *UpmsHandler) Create(ctx *gin.Context, req *model.Upms) error {
	if err := l.db.WithContext(ctx).Create(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("权限创建失败, err: %v", err))
		return fmt.Errorf("权限创建失败")
	}
	return nil
}
func (l *UpmsHandler) Put(ctx *gin.Context, id types.SearchId, req *model.Upms) (*model.Upms, error) {
	if err := l.db.WithContext(ctx).Model(&model.Upms{}).Where("id = ?", id.Id).Updates(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("权限更新失败, err: %v", err))
		return nil, fmt.Errorf("权限更新失败")
	}
	// 查询出最新的记录
	if err := l.db.WithContext(ctx).Model(&model.Upms{}).Where("id = ?", id.Id).First(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("权限更新失败, err: %v", err))
		return nil, fmt.Errorf("权限更新失败")
	}
	return req, nil
}
func (l *UpmsHandler) Delete(ctx *gin.Context, req types.SearchId) error {
	if err := l.db.WithContext(ctx).Model(&model.Upms{}).Where("id = ?", req.Id).Delete(&model.Upms{}).Error; err != nil {
		l.l.Error(fmt.Sprintf("权限删除失败, err: %v", err))
		return fmt.Errorf("权限删除失败")
	}
	return nil
}

func NewUpmsHandler() *UpmsHandler {
	return &UpmsHandler{
		l:  global.L.Named(apps.AppName).Named(apps.AppUpms).Named("logic"),
		db: global.DB.GetDb(),
	}
}

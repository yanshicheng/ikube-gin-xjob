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

var _ service.MenuService = (*MenuLogic)(nil)

type MenuLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (l *MenuLogic) List(ctx *gin.Context, req types2.MenuSearchReq) ([]*model.Menu, error) {

	return nil, nil
}

func (l *MenuLogic) Create(ctx *gin.Context, req *model.Menu) error {
	// 判断 ParentId 是否存在
	if *req.ParentId != 0 {
		var parent model.Menu
		if err := l.db.WithContext(ctx).Model(&model.Menu{}).Where("id = ?", req.ParentId).First(&parent).Error; err != nil {
			l.l.Error(fmt.Sprintf("菜单创建失败, err: %v", err))
			return fmt.Errorf("菜单创建失败，parentId 不存在")
		}
	}

	if err := l.db.WithContext(ctx).Create(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("菜单创建失败, err: %v", err))
		return fmt.Errorf("菜单创建失败")
	}
	return nil
}
func (l *MenuLogic) Put(ctx *gin.Context, id types.SearchId, req *model.Menu) (*model.Menu, error) {
	if err := l.db.WithContext(ctx).Model(&model.Menu{}).Where("id = ?", id.Id).Updates(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("菜单更新失败, err: %v", err))
		return nil, fmt.Errorf("菜单更新失败")
	}
	// 查询出最新的记录
	if err := l.db.WithContext(ctx).Model(&model.Menu{}).Where("id = ?", id.Id).First(req).Error; err != nil {
		l.l.Error(fmt.Sprintf("菜单更新失败, err: %v", err))
		return nil, fmt.Errorf("菜单更新失败")
	}
	return req, nil
}
func (l *MenuLogic) Delete(ctx *gin.Context, req types.SearchId) error {
	if err := l.db.WithContext(ctx).Model(&model.Menu{}).Where("id = ?", req.Id).Delete(&model.Menu{}).Error; err != nil {
		l.l.Error(fmt.Sprintf("菜单删除失败, err: %v", err))
		return fmt.Errorf("菜单删除失败")
	}
	return nil
}

func NewMenuLogic() *MenuLogic {
	return &MenuLogic{
		l:  global.L.Named(apps.AppName).Named(apps.AppMenu).Named("logic"),
		db: global.DB.GetDb(),
	}
}

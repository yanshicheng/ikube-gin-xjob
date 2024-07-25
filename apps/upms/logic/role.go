package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/upms"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/service"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/sql"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ service.RoleService = (*RoleLogic)(nil)

type RoleLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (r *RoleLogic) List(c *gin.Context, search types2.RoleSearchReq) (*types.QueryResponse, error) {
	var list []*model.Role
	db := r.db.WithContext(c).Model(&model.Role{})
	// 设置排序
	db = db.Order(fmt.Sprintf("%s %s", "ID", search.Sort))
	// 模糊查询
	if search.Name != "" {
		db = db.Where("name like ?", "%"+search.Name+"%")
	}
	// 打印 sql 语句
	queryRes, err := sql.GetQueryResponse(db, search.Pagination, list)
	if err != nil {
		r.l.Error(fmt.Sprintf("查询角色失败: %s", err.Error()))
		return nil, fmt.Errorf("查询角色失败")
	}
	return queryRes, nil
}
func (r *RoleLogic) Create(c *gin.Context, req *model.Role) error {
	if err := r.db.WithContext(c).Create(req).Error; err != nil {
		r.l.Error(fmt.Sprintf("创建角色失败: %s", err.Error()))
		return fmt.Errorf("创建角色失败")
	}
	return nil
}
func (r *RoleLogic) Put(c *gin.Context, search types.SearchId, req *types2.RoleUpdateRequest) (*model.Role, error) {
	// 只允许修改名称
	if err := r.db.WithContext(c).Model(&model.Role{}).Where("id = ?", search.Id).Updates(model.Role{Name: req.Name}).Error; err != nil {
		r.l.Error(fmt.Sprintf("修改角色失败: %s", err.Error()))
		return nil, fmt.Errorf("修改角色失败")
	}
	updatedRole := model.Role{}
	if err := r.db.WithContext(c).Model(&model.Role{}).Where("id = ?", search.Id).First(&updatedRole).Error; err != nil {
		r.l.Error(fmt.Sprintf("查询角色失败: %s", err.Error()))
		return nil, fmt.Errorf("查询角色失败")
	}
	return &updatedRole, nil
}

func (r *RoleLogic) Delete(c *gin.Context, id types.SearchId) error {
	// 检查 ikubexjob_user_account 表中是否存在这个 role_id
	var count int64
	err := r.db.Table("ikubexjob_user_account").Where("role_id = ?", id.Id).Count(&count).Error
	if err != nil {
		r.l.Error(fmt.Sprintf("查询角色失败: %s", err.Error()))
		return fmt.Errorf("查询角色失败")
	}

	// 如果存在引用，则不删除并返回错误
	if count > 0 {
		r.l.Error(fmt.Sprintf("无法删除角色，因为它在 ikubexjob_user_account 表中仍有引用"))
		return fmt.Errorf("无法删除角色，角色正在使用中")
	}
	result := r.db.WithContext(c).Model(&model.Role{}).Where("id = ?", id.Id).Delete(&model.Role{})
	if err := result.Error; err != nil {
		r.l.Error(fmt.Sprintf("删除角色失败: %s", err.Error()))
		return fmt.Errorf("删除角色失败")
	}
	// 检查是否实际上删除了记录
	if result.RowsAffected == 0 {
		r.l.Error("删除角色失败: 未找到指定的角色")
		return fmt.Errorf("删除角色失败: 未找到指定的角色")
	}
	return nil
}

// Config 只需要保证 全局对象Config和全局Logger已经加载完成
func (r *RoleLogic) Config() {
	r.l = global.L.Named(apps.AppName).Named(apps.AppRole).Named("logic")
	r.db = global.DB.GetDb()
}

func (r *RoleLogic) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppRole)
}

func NewRoleLogic() *RoleLogic {
	return &RoleLogic{
		l:  global.L.Named(apps.AppName).Named(apps.AppRole).Named("logic"),
		db: global.DB.GetDb(),
	}
}

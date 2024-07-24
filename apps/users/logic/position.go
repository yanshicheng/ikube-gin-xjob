package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/service"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ service.PositionService = (*PositionLogic)(nil)

var positionLogic = &PositionLogic{}

type PositionLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (o *PositionLogic) List(c *gin.Context, search types2.PositionListSearchReq) ([]*model.Position, error) {
	var list []*model.Position
	if search.OrganizationId != 0 {
		if err := o.db.WithContext(c).Model(&model.Position{}).Where("organization_id = ?", search.OrganizationId).Find(&list).Error; err != nil {
			o.l.Error(fmt.Sprintf("获取机构职位信息失败: %s， ID: %d", err.Error(), search.OrganizationId))
			return nil, fmt.Errorf("获取机构职位信息失败,ID: %d", search.OrganizationId)
		}
	} else {
		if err := o.db.WithContext(c).Model(&model.Position{}).Find(&list).Error; err != nil {
			o.l.Error(fmt.Sprintf("获取职位所有信息失败: %s", err.Error()))
			return nil, fmt.Errorf("获取职位所有信息失败,查询失败")
		}
	}
	return list, nil
}
func (o *PositionLogic) Create(c *gin.Context, req *model.Position) error {
	// 创建职位，先查询机构是否存在
	var org model.Organization
	if err := o.db.WithContext(c).Model(&model.Organization{}).Where("id = ?", req.OrganizationId).First(&org).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构失败无法创建职位: %s， ID: %d", err.Error(), req.OrganizationId))
		return fmt.Errorf("查询机构失败无法创建职位,ID: %d", req.OrganizationId)
	}
	// 创建机构，只允许创建主体机构
	if org.Level == 0 {
		return fmt.Errorf("创建职位失败，组织不是主体")
	}
	if err := o.db.WithContext(c).Model(&model.Position{}).Create(&req).Error; err != nil {
		o.l.Error(fmt.Sprintf("创建职位失败: %s", err.Error()))
		return fmt.Errorf("职位创建失败")
	}
	return nil
}
func (o *PositionLogic) Put(c *gin.Context, search types.SearchId, req *model.Position) (*model.Position, error) {
	// 只允许修改名称
	//updates := map[string]interface{}{"name": pos.Name}
	if err := o.db.WithContext(c).Model(&model.Position{}).Where("id = ?", search.Id).Updates(map[string]string{"name": req.Name}).Error; err != nil {
		o.l.Error(fmt.Sprintf("更新职位失败: %s", err.Error()))
		return nil, fmt.Errorf("更新职位失败: %d", search.Id)
	}

	// 重新获取更新后的职位信息
	var updatedPosition model.Position
	if err := o.db.WithContext(c).Model(&model.Position{}).Where("id = ?", search.Id).First(&updatedPosition).Error; err != nil {
		o.l.Error(fmt.Sprintf("获取更新后的职位信息失败: %s", err.Error()))
		return nil, fmt.Errorf("获取更新后的职位信息失败")
	}
	return &updatedPosition, nil
}

func (o *PositionLogic) Delete(c *gin.Context, id types.SearchId) error {
	// 先判断有没有用户占用职位
	var count int64
	if err := o.db.WithContext(c).Model(&model.Account{}).Where("position_id = ?", id.Id).Count(&count).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询职位被占用失败: %s", err.Error()))
		return fmt.Errorf("查询职位被占用失败")
	}
	if count > 0 {
		o.l.Error(fmt.Sprintf("职位被占用无法删除: %d", id.Id))
		return fmt.Errorf("职位被占用无法删除: %d", id.Id)
	}
	// 执行删除操作，并获取结果
	result := o.db.WithContext(c).Model(&model.Position{}).Where("id = ?", id.Id).Delete(&model.Position{})

	// 检查是否有错误发生
	if err := result.Error; err != nil {
		o.l.Error(fmt.Sprintf("删除职位失败: %s", err.Error()))
		return fmt.Errorf("删除职位失败: %s", err.Error())
	}

	// 检查是否实际上删除了记录
	if result.RowsAffected == 0 {
		o.l.Error("删除职位失败: 未找到指定的职位")
		return fmt.Errorf("删除职位失败: 未找到指定的职位")
	}

	return nil
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (o *PositionLogic) Config() {
	o.l = global.L.Named(users.AppName).Named(users.AppPosition).Named("logic")
	o.db = global.DB.GetDb()
}

func (o *PositionLogic) Name() string {
	return fmt.Sprintf("%s.%s", users.AppName, users.AppPosition)
}

func init() {
	// 注册
	router.RegistryLogic(positionLogic)
}

package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/service"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 接口检查
var _ service.OrganizationService = (*OrganizationLogic)(nil)

//var logic = &OrganizationLogic{}

type OrganizationLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (o *OrganizationLogic) Get(c *gin.Context, search types2.OrganizationGetSearchReq) ([]*model.Organization, error) {
	o.l.Info(fmt.Sprintf("查询机构信息, name: %s", search.Name))
	var orsSlice []model.Organization
	if err := o.db.WithContext(c).Where("name like ?", search.Name+"%").Find(&orsSlice).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构信息失败, id: %s, error: %s", search.Name, err.Error()))
		return nil, err
	}
	resultOrgSlice := make([]*model.Organization, 0)
	for _, org := range orsSlice {
		resultOrg, err := o.buildSingleParent(&org)
		if err != nil {
			o.l.Error(fmt.Sprintf("构建父级层级结构失败, id: %d, error: %s", org.ID, err.Error()))
			return nil, err
		}
		resultOrgSlice = append(resultOrgSlice, resultOrg)
	}
	return resultOrgSlice, nil
}

// 递归获某个节点的所有父节点
func (o *OrganizationLogic) buildSingleParent(org *model.Organization) (*model.Organization, error) {
	if org.ParentId == 0 {
		return org, nil // 如果没有父级，返回当前组织
	}

	var parent model.Organization
	if err := o.db.Where("id = ?", org.ParentId).First(&parent).Error; err != nil {
		return nil, err
	}
	// 将当前组织添加到父级的 Children 中
	parent.Children = append(parent.Children, org)
	// 判断是不是根节点
	if parent.ParentId == 0 {
		return &parent, nil
	}
	// 递归构建父级结构
	parentOrg, err := o.buildSingleParent(&parent)
	if err != nil {
		return nil, err
	}

	return parentOrg, nil
}

// 递归获某个节点的所有父节点
func (o *OrganizationLogic) buildParentHierarchy(org *model.Organization, orgSlice []*model.Organization) ([]*model.Organization, error) {
	// 先判断是否存在列表中
	args, ok := o.recursiveAdd(org, orgSlice)
	if ok {
		return args, nil
	}

	var parent model.Organization
	if err := o.db.Where("id = ?", org.ParentId).First(&parent).Error; err != nil {
		return nil, err
	}
	// 将当前组织添加到父级的 Children 中
	parent.Children = append(parent.Children, org)
	args, ok = o.recursiveAdd(&parent, orgSlice)
	if ok {
		return args, nil
	} else {
		orgSlice, err := o.buildParentHierarchy(&parent, orgSlice)
		if err != nil {
			return nil, err
		}
		return orgSlice, nil
	}
}

// 递归添加节点
func (o *OrganizationLogic) recursiveAdd(org *model.Organization, orgSlice []*model.Organization) ([]*model.Organization, bool) {
	// 判断当前是否为主节点
	if org.ParentId == 0 {
		for _, og := range orgSlice {
			if og.ID == org.ID {
				// 如果已经存在，则跳过
				return orgSlice, true
			}
		}
		// 如果不存在，则添加
		orgSlice = append(orgSlice, org)
		return orgSlice, true
	}
	// 当前节点为子节点
	for _, og := range orgSlice {
		if og.ID == org.ParentId {
			// 如果已经存在，则跳过
			og.Children = append(og.Children, org)
			return orgSlice, true
		}
		// 如果当前节点存在子节点，则递归添加
		if og.Children != nil {
			o.recursiveAdd(org, og.Children)
		}
	}
	// 如果都不满足，则返回false
	return orgSlice, false
}
func (o *OrganizationLogic) List(c *gin.Context, search types2.OrganizationGetSearchReq) ([]*model.Organization, error) {
	if search.Name == "" {
		var allOrgSlice []model.Organization
		if err := o.db.WithContext(c).Find(&allOrgSlice).Error; err != nil {
			o.l.Error(fmt.Sprintf("查询机构信息失败, error: %s", err.Error()))
			return nil, fmt.Errorf("查询机构信息失败")
		}
		orgMap := make(map[uint]*model.Organization, len(allOrgSlice))
		for i := range allOrgSlice {
			orgMap[allOrgSlice[i].ID] = &allOrgSlice[i]
		}
		// 构建树形结构
		var orgTree []*model.Organization
		for i := range allOrgSlice {
			org := &allOrgSlice[i]
			if org.ParentId != 0 {
				if parent, ok := orgMap[org.ParentId]; ok {
					if parent.Children == nil {
						parent.Children = make([]*model.Organization, 0)
					}
					parent.Children = append(parent.Children, org)
				}
			} else {
				orgTree = append(orgTree, org)
			}
		}
		return orgTree, nil
	} else {
		var orgSlice []model.Organization
		if err := o.db.WithContext(c).Where("name like ?", search.Name+"%").Find(&orgSlice).Error; err != nil {
			o.l.Error(fmt.Sprintf("查询机构信息失败, id: %s, error: %s", search.Name, err.Error()))
			return nil, err
		}
		o.l.Debug(fmt.Sprintf("查询机构信息: %+v, %d", orgSlice, len(orgSlice)))
		var resultOrgSlice []*model.Organization
		for _, org := range orgSlice {
			o.l.Debug(fmt.Sprintf("查询机构信息org: %+v", org))
			resOrgSlice, err := o.buildParentHierarchy(&org, resultOrgSlice)
			if err != nil {
				o.l.Error(fmt.Sprintf("构建父级层级结构失败, id: %d, error: %s", org.ID, err.Error()))
				return nil, err
			}
			resultOrgSlice = make([]*model.Organization, len(resOrgSlice))
			copy(resultOrgSlice, resOrgSlice)
		}
		o.l.Debug(fmt.Sprintf("查询机构信息resultOrg: %+v", resultOrgSlice))
		return resultOrgSlice, nil
	}
}

func (o *OrganizationLogic) Put(c *gin.Context, id types.SearchId, org *model.Organization) (*model.Organization, error) {
	// 先查询出来
	var oldOrg model.Organization
	if err := o.db.WithContext(c).Where("id = ?", id.Id).First(&oldOrg).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return nil, fmt.Errorf("查询机构信息失败")
	}
	// 判断是否为主节点，如果是主节点则不允许修改 ParentId
	if oldOrg.ParentId == 0 {
		if org.ParentId != 0 {
			o.l.Error(fmt.Sprintf("主节点不允许修改 ParentId, id: %d", id.Id))
			return nil, fmt.Errorf("主节点不允许修改 ParentId")
		}
	}
	// 判断是否存在子节点
	var childOrgSlice []model.Organization
	if err := o.db.WithContext(c).Where("parent_id = ?", id.Id).Find(&childOrgSlice).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询子节点信息失败, id: %d, error: %s", id.Id, err.Error()))
		return nil, fmt.Errorf("查询子节点信息失败")
	}
	if len(childOrgSlice) > 0 {
		// 如果存在子节点，则不允许修改 ParentId
		if org.ParentId != oldOrg.ParentId {
			o.l.Error(fmt.Sprintf("存在子节点不允许修改 ParentId, id: %d", id.Id))
			return nil, fmt.Errorf("存在子节点不允许修改 ParentId")
		}
	}
	// 修改操作
	oldOrg.Desc = org.Desc
	oldOrg.Name = org.Name
	oldOrg.ParentId = org.ParentId
	// 保存
	if err := o.db.WithContext(c).Save(&oldOrg).Error; err != nil {
		o.l.Error(fmt.Sprintf("更新机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return nil, fmt.Errorf("更新机构信息失败")
	}
	return &oldOrg, nil
}

func (o *OrganizationLogic) Create(c *gin.Context, req *model.Organization) error {
	if err := o.db.WithContext(c).Create(req).Error; err != nil {
		o.l.Error(fmt.Sprintf("创建机构信息失败, error: %s", err.Error()))
		return err
	}
	return nil
}

func (o *OrganizationLogic) Delete(c *gin.Context, id types.SearchId) error {
	// 删除机构首先查询出来
	var org model.Organization
	if err := o.db.WithContext(c).Where("id = ?", id.Id).First(&org).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return fmt.Errorf("查询机构信息失败")
	}
	// 判断是否有子节点
	var childOrgSlice []model.Organization
	if err := o.db.WithContext(c).Where("parent_id = ?", id.Id).Find(&childOrgSlice).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return fmt.Errorf("查询机构信息失败")
	}
	if len(childOrgSlice) > 0 {
		o.l.Error(fmt.Sprintf("机构下存在子节点, id: %d", id.Id))
		return fmt.Errorf("机构下存在子节点")
	}
	// 查询是否有职位
	var positionSlice []model.Position
	if err := o.db.WithContext(c).Where("organization_id = ?", id.Id).Find(&positionSlice).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询职位信息失败, id: %d, error: %s", id.Id, err.Error()))
		return fmt.Errorf("查询职位信息失败")
	}
	if len(positionSlice) > 0 {
		o.l.Error(fmt.Sprintf("机构下存在职位, id: %d", id.Id))
		return fmt.Errorf("机构下存在职位")
	}
	// 删除
	result := o.db.WithContext(c).Where("id = ?", id.Id).Delete(&model.Organization{})
	if result.Error != nil {
		o.l.Error(fmt.Sprintf("删除机构信息失败, id: %d, error: %s", id.Id, result.Error.Error()))
		return fmt.Errorf("删除机构信息失败")
	}
	if result.RowsAffected == 0 {
		o.l.Error(fmt.Sprintf("删除机构信息失败, id: %d, error: %s", id.Id, result.Error.Error()))
		return fmt.Errorf("删除机构信息失败,职位不存在")
	}
	return nil
}

// Config 只需要保证 全局对象Config和全局Logger已经加载完成
func (o *OrganizationLogic) Config() {
	o.l = global.L.Named(apps.AppName).Named(apps.AppOrganization).Named("logic")
	o.db = global.DB.GetDb()
}

func (o *OrganizationLogic) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppOrganization)
}

//func init() {
//	// 注册
//	router.RegistryLogic(logic)
//}

func NewOrganizationLogic() *OrganizationLogic {
	return &OrganizationLogic{
		l:  global.L.Named(apps.AppName).Named(apps.AppOrganization).Named("logic"),
		db: global.DB.GetDb(),
	}
}

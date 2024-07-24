package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/service"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/errorx"
	"github.com/yanshicheng/ikube-gin-xjob/common/sql"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"github.com/yanshicheng/ikube-gin-xjob/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var _ service.AccountService = (*AccountLogic)(nil)

var accountLogic = &AccountLogic{}

type AccountLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (l *AccountLogic) Get(c *gin.Context, search types.SearchId) (*model.Account, error) {
	var account *model.Account
	// 查询账号详情
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", search.Id).First(&account).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询账号详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询账号详情失败")
	}
	var position *model.Position
	if err := l.db.WithContext(c).Model(&model.Position{}).Where("id = ?", account.PositionId).First(&position).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询职位详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询职位详情失败")
	}
	account.PositionName = position.Name
	// 查找部门信息
	var organization *model.Organization
	if err := l.db.WithContext(c).Model(&model.Organization{}).Where("id = ?", account.OrganizationId).First(&organization).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询部门详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询部门详情失败")
	}
	account.OrganizationName = organization.Name
	OrganizationTreeName, err := organization.GetFullHierarchy(l.db)
	if err != nil {
		l.l.Error(fmt.Sprintf("查询部门详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询部门详情失败")
	}
	account.OrganizationTreeName = OrganizationTreeName
	// 补充职位信息，
	return account, nil
}
func (l *AccountLogic) List(c *gin.Context, query types2.AccountQueryReq) (*types.QueryResponse, error) {
	var list []*model.Account
	// 设置排序
	db := l.db.WithContext(c).Model(&model.Account{})
	db = db.Order(fmt.Sprintf("%s %s", "ID", query.Sort))
	// 模糊查询
	if query.Account != "" {
		db = db.Where("title like ?", "%d"+query.Account+"%")
	}
	if query.UserName != "" {
		db = db.Where("title like ?", "%d"+query.UserName+"%")
	}
	if query.Mobile != "" {
		db = db.Where("title like ?", "%d"+query.Mobile+"%")
	}
	if query.WorkNumber != "" {
		db = db.Where("title like ?", "%d"+query.WorkNumber+"%")
	}
	if query.Email != "" {
		db = db.Where("title like ?", "%d"+query.Email+"%")
	}
	if !query.IsDisabled {
		db = db.Where("is_disabled = ?", query.IsDisabled)
	}
	if !query.IsFrozen {
		db = db.Where("is_frozen = ?", query.IsFrozen)
	}
	if !query.IsLeave {
		db = db.Where("is_leave = ?", query.IsLeave)
	}
	if query.PositionId != nil {
		db = db.Where("position_id = ?", query.PositionId)
	}
	if query.OrganizationId != nil {
		db = db.Where("organization_id = ?", query.OrganizationId)
	}
	queryRes, err := sql.GetQueryResponse(db, query.Pagination, list)
	if err != nil {
		l.l.Error(fmt.Sprintf("查询数据失败: %s", err.Error()))
		return nil, fmt.Errorf("查询数据失败")
	}
	return queryRes, nil
}

func (l *AccountLogic) Create(c *gin.Context, data *types2.AccountCreateReq) (*model.Account, error) {
	// 创建账号，先查询机构是否存在
	account := model.Account{
		Account:        data.Account,
		UserName:       data.UserName,
		Email:          data.Email,
		HireDate:       data.HireDate,
		WorkNumber:     data.WorkNumber,
		Mobile:         data.Mobile,
		OrganizationId: data.OrganizationId,
		PositionId:     data.PositionId,
	}
	// 判断组织是否存在
	if err := l.db.WithContext(c).Model(&model.Organization{}).Where("id = ?", data.OrganizationId).First(&model.Organization{}).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询部门详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询部门详情失败")
	}
	// 判断职位是否存在
	if err := l.db.WithContext(c).Model(&model.Position{}).Where("id = ?", data.PositionId).First(&model.Position{}).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询职位详情失败: %s", err.Error()))
		return nil, fmt.Errorf("查询职位详情失败")
	}
	// 创建账号， 进行密码加密
	err := account.SetPassword(utils.GeneratePassword())
	if err != nil {
		return nil, err
	}
	// 设置默认头像
	account.Icon = utils.GenerateIcon()
	// 设置必须重置密码
	account.IsChangePassword = true
	// 创建账号
	if err := l.db.WithContext(c).Create(&account).Error; err != nil {
		l.l.Error(fmt.Sprintf("创建账号失败: %s", err.Error()))
		return nil, fmt.Errorf("创建账号失败")
	}
	return &account, nil
}
func (l *AccountLogic) Put(c *gin.Context, search types.SearchId, new *types2.AccountCreateReq) (*model.Account, error) {
	// 只允许修改名称
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", search.Id).Updates(new).Error; err != nil {
		l.l.Error(fmt.Sprintf("更新账号失败: %s", err.Error()))
		return nil, fmt.Errorf("更新账号失败: %d", search.Id)
	}

	// 重新获取更新后的账号信息
	var updatedPosition model.Account
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", search.Id).First(&updatedPosition).Error; err != nil {
		l.l.Error(fmt.Sprintf("获取更新后的账号信息失败: %s", err.Error()))
		return nil, fmt.Errorf("获取更新后的账号信息失败")
	}
	return &updatedPosition, nil
}

func (l *AccountLogic) Delete(c *gin.Context, id types.SearchId) error {
	// 执行删除操作，并获取结果
	result := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", id.Id).Delete(&model.Account{})

	// 检查是否有错误发生
	if err := result.Error; err != nil {
		l.l.Error(fmt.Sprintf("删除账号失败: %s", err.Error()))
		return fmt.Errorf("删除账号失败: %s", err.Error())
	}

	// 检查是否实际上删除了记录
	if result.RowsAffected == 0 {
		l.l.Error("删除账号失败: 未找到指定的账号")
		return fmt.Errorf("删除账号失败: 未找到指定的账号")
	}
	return nil
}

func (l *AccountLogic) ChangePassword(c *gin.Context, req *types2.AccountChangePasswordReq) error {
	var account model.Account
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("account = ?", req.Account).First(&account).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询账号失败: %s", err.Error()))
		return fmt.Errorf("查询账号失败")
	}
	// 对原始密码解密
	oldPassword, err := utils.DecodeBase64Password(req.Password)
	if err != nil {
		l.l.Error(fmt.Sprintf("解密密码失败: %s", err.Error()))
		return fmt.Errorf("解密密码失败")
	}
	// 验证密码
	if !account.CheckPassword(oldPassword) {
		l.l.Error(fmt.Sprintf("用户:%s  密码错误: %s", req.Account, oldPassword))
		return fmt.Errorf("用户名或密码错误")
	}
	// 判断用户是否离职 或者 是否为禁用，
	if account.IsLeave || account.IsDisabled {
		l.l.Error(fmt.Sprintf("用户:%s  已被禁用", req.Account))
		return fmt.Errorf("用户已被禁用，请联系管理员")
	}
	// 解析新密码
	newPassword, err := utils.DecodeBase64Password(req.NewPassword)
	if err != nil {
		l.l.Error(fmt.Sprintf("解密密码失败: %s", err.Error()))
		return fmt.Errorf("解密密码失败")
	}
	// 验证密码复杂度要求
	if !utils.CheckPasswordComplexity(newPassword) {
		l.l.Error(fmt.Sprintf("密码复杂度要求:密码: %s", newPassword))
		return fmt.Errorf("不满足密码复杂度要求，要求: 至少 12 位包含数字、字母、特殊字符")
	}
	// 判断新密码是否与旧密码相同
	if oldPassword == newPassword {
		l.l.Error(fmt.Sprintf("新密码不能与旧密码相同: %s", newPassword))
		return fmt.Errorf("新密码不能与旧密码相同")
	}
	// 验证密码
	if !account.CheckPassword(newPassword) {
		l.l.Error(fmt.Sprintf("用户名或者密码错误: %s", newPassword))
		return fmt.Errorf("用户名或者密码错误")
	}
	if err := account.SetPassword(newPassword); err != nil {
		l.l.Error(fmt.Sprintf("设置密码失败: %s", err.Error()))
		return fmt.Errorf("设置密码失败")
	}
	// 设置密码重置标识
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", account.ID).Update("is_reset_password", 0).Error; err != nil {
		l.l.Error(fmt.Sprintf("设置密码重置标识失败: %s", err.Error()))
		return fmt.Errorf("设置密码重置标识失败")
	}
	// 设置新密码
	return nil

}

func (l *AccountLogic) RestPassword(c *gin.Context, req *types2.AccountRestPasswordReq) error {
	var account model.Account
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("account = ?", req.Account).First(&account).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询账号失败: %s", err.Error()))
		return fmt.Errorf("查询账号失败")
	}
	// 创建账号， 进行密码加密
	err := account.SetPassword(utils.GeneratePassword())
	// 设置必须修改密码
	account.IsChangePassword = true
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", account.ID).Updates(account).Error; err != nil {
		l.l.Error(fmt.Sprintf("更新账号失败: %s", err.Error()))
		return fmt.Errorf("更新账号失败")
	}
	if err != nil {
		return err
	}
	return nil
}
func (l *AccountLogic) Login(c *gin.Context, req *types2.AccountLoginReq) (*utils.JWTResponse, errorx.ErrorCode, error) {
	var account model.Account
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("account = ?", req.Account).First(&account).Error; err != nil {
		l.l.Error(fmt.Sprintf("查询账号失败: %s", err.Error()))
		return nil, errorx.ErrGeneric, fmt.Errorf("用户名或密码错误")
	}
	// 密码解密
	password, err := utils.DecodeBase64Password(req.Password)
	if err != nil {
		l.l.Error(fmt.Sprintf("解密密码失败: %s", err.Error()))
		return nil, errorx.ErrGeneric, fmt.Errorf("解密密码失败")
	}
	if !account.CheckPassword(password) {
		l.l.Error(fmt.Sprintf("用户:%s  密码错误: %s", req.Account, password))
		return nil, errorx.ErrGeneric, fmt.Errorf("用户名或密码错误")
	}
	if account.IsDisabled {
		l.l.Error(fmt.Sprintf("用户:%s  已被禁用", req.Account))
		return nil, errorx.ErrGeneric, fmt.Errorf("用户已被禁用，请联系管理员")
	}
	if account.IsLeave {
		l.l.Error(fmt.Sprintf("用户:%s  已被离职", req.Account))
		return nil, errorx.ErrGeneric, fmt.Errorf("用户已被离职，请联系管理员")
	}
	if account.IsChangePassword {
		l.l.Error(fmt.Sprintf("用户:%s  需要重置密码", req.Account))
		return nil, errorx.ErrNeedResetPassword, fmt.Errorf("用户需要重置密码，请联系管理员")
	}
	if err := l.db.WithContext(c).Model(&model.Account{}).Where("id = ?", account.ID).Update("last_login_time", time.Now()).Error; err != nil {
		l.l.Error(fmt.Sprintf("更新账号登录时间失败: %s", err.Error()))
		return nil, errorx.ErrGeneric, fmt.Errorf("更新账号登录时间失败")
	}

	// var
	var applicationRole utils.ApplicationRole
	fmt.Println(applicationRole)
	return nil, errorx.ErrGeneric, nil
}
func (l *AccountLogic) Logout(*gin.Context) error {
	return nil
}
func (l *AccountLogic) ChangeIcon(*gin.Context) (types2.AccountIconResp, error) {
	return types2.AccountIconResp{}, nil
}

// Config 只需要保证 全局对象Config和全局Logger已经加载完成
func (l *AccountLogic) Config() {
	l.l = global.L.Named(apps.AppName).Named(apps.AppAccount).Named("logic")
	l.db = global.DB.GetDb()
}

func (l *AccountLogic) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppAccount)
}

func init() {
	// 注册
	router.RegistryLogic(accountLogic)
}

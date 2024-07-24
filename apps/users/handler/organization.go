package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/logic"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/response"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"go.uber.org/zap"
)

var _ router.GinService = (*OrganizationHandler)(nil)
var organizationHandler = &OrganizationHandler{}

type OrganizationHandler struct {
	l   *zap.Logger
	svc *logic.OrganizationLogic
}

func (h *OrganizationHandler) PublicRegistry(r gin.IRouter) {

}

// AuthRegistry 注册认证接口
func (h *OrganizationHandler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group(fmt.Sprintf("%s/%s", apps.AppName, apps.AppOrganization))
	{
		group.GET("/", h.list)
		group.GET("/:id", h.get)
		group.POST("/", h.create)
		group.PUT("/:id", h.put)
		group.DELETE("/:id", h.delete)
	}
}
func (h *OrganizationHandler) get(c *gin.Context) {
	var search types2.OrganizationGetSearchReq
	if err := c.ShouldBindUri(&search); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	h.l.Info(fmt.Sprintf("查询参数: %+v", search))
	if s, err := h.svc.Get(c, search); err != nil {
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("查询成功: %+v", s))
		response.SuccessSlice(c, s)
	}

}
func (h *OrganizationHandler) list(c *gin.Context) {
	var search types2.OrganizationGetSearchReq
	if err := c.ShouldBindQuery(&search); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("查询参数: %v", search))

	if s, err := h.svc.List(c, search); err != nil {
		h.l.Error(fmt.Sprintf("数据查询失败: %s", err))
		response.FailedStr(c, err.Error())
	} else {
		response.SuccessSlice(c, s)
	}
}

func (h *OrganizationHandler) create(c *gin.Context) {
	var org model.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.Create(c, &org); err != nil {
		h.l.Error(fmt.Sprintf("数据创建失败: %s", err))
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("创建成功: %+v", org))
		response.SuccessMap(c, org)
	}

}
func (h *OrganizationHandler) put(c *gin.Context) {
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	var org model.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	if newOrg, err := h.svc.Put(c, id, &org); err != nil {
		h.l.Error(fmt.Sprintf("数据更新失败: %s", err))
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("更新成功: %+v", org))
		response.SuccessMap(c, newOrg)
	}

}
func (h *OrganizationHandler) delete(c *gin.Context) {
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}

	if err := h.svc.Delete(c, id); err != nil {
		h.l.Error(fmt.Sprintf("数据删除失败: %s", err))
		response.FailedStr(c, err.Error())
		return
	}
	h.l.Debug(fmt.Sprintf("删除成功: %+v", id))
	response.SuccessMap(c, nil)
}

func (h *OrganizationHandler) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppOrganization)
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *OrganizationHandler) Config() {
	h.l = global.L.Named(apps.AppName).Named(apps.AppOrganization).Named("handler")
	h.svc = logic.NewOrganizationLogic()
}

func init() {
	router.RegistryGinRouter(organizationHandler)
}

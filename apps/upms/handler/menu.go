package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/upms"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/logic"
	"github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/upms/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/response"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"go.uber.org/zap"
)

var _ router.GinService = (*MenuHandler)(nil)
var menuHandler = &MenuHandler{}

type MenuHandler struct {
	l   *zap.Logger
	svc *logic.MenuLogic
}

func (h *MenuHandler) PublicRegistry(gin.IRouter) {

}

// AuthRegistry 注册认证接口
func (h *MenuHandler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group(fmt.Sprintf("%s/%s", apps.AppName, apps.AppMenu))
	{
		group.GET("/", h.list)
		group.POST("/", h.create)
		group.PUT("/:id", h.put)
		group.DELETE("/:id", h.delete)
	}
}

func (h *MenuHandler) list(c *gin.Context) {
	search := types2.MenuSearchReq{}
	if err := c.ShouldBindQuery(&search); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("查询参数: %+v", search))
	list, err := h.svc.List(c, search)
	if err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessSlice(c, list)
}

func (h *MenuHandler) create(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.Create(c, &req); err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessMap(c, req)
}

func (h *MenuHandler) put(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("修改参数: %+v, 修改id: %d", req, id))
	if newPosition, err := h.svc.Put(c, id, &req); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessMap(c, newPosition)
	}

}

func (h *MenuHandler) delete(c *gin.Context) {
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("删除id: %d", id))
	if err := h.svc.Delete(c, id); err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessMap(c, nil)
}

func (h *MenuHandler) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppMenu)
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *MenuHandler) Config() {
	h.l = global.L.Named(apps.AppName).Named(apps.AppMenu).Named("handler")
	h.svc = logic.NewMenuLogic()
}

func init() {
	router.RegistryGinRouter(menuHandler)
}

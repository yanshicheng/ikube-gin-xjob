package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apps "github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/logic"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/response"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"go.uber.org/zap"
)

var _ router.GinService = (*AccountHandler)(nil)
var accountHandler = &AccountHandler{}

type AccountHandler struct {
	l   *zap.Logger
	svc *logic.AccountLogic
}

func (h *AccountHandler) PublicRegistry(r gin.IRouter) {
	group := r.Group(fmt.Sprintf("%s/%s", apps.AppName, apps.AppAccount))
	group.POST("/login", h.login)
	group.POST("/logout", h.login)
	// 重置密码接口
	group.POST("/changePassword", h.changePassword)

}

// AuthRegistry 注册认证接口
func (h *AccountHandler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group(fmt.Sprintf("%s/%s", apps.AppName, apps.AppAccount))
	{
		group.GET("/", h.list)
		group.GET("/:id", h.get)
		group.POST("/", h.create)
		group.PUT("/:id", h.put)
		group.DELETE("/:id", h.delete)
		group.POST("resetPassword", h.resetPassword)
	}

}

func (h *AccountHandler) changePassword(c *gin.Context) {
	var req types2.AccountChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.ChangePassword(c, &req); err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessStr(c, "密码修改成功!")
}
func (h *AccountHandler) resetPassword(c *gin.Context) {
	var req types2.AccountRestPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.RestPassword(c, &req); err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessStr(c, "密码重置成功!")
}

func (h *AccountHandler) login(c *gin.Context) {
	var req types2.AccountLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if token, errCode, err := h.svc.Login(c, &req); err != nil {
		response.FailedCode(c, errCode, err.Error())
		return
	} else {
		response.SuccessMap(c, token)
	}
}

func (h *AccountHandler) get(c *gin.Context) {
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if account, err := h.svc.Get(c, id); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessMap(c, account)
	}

}

func (h *AccountHandler) list(c *gin.Context) {
	var search types2.AccountQueryReq
	if err := c.ShouldBindQuery(&search); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if accounts, err := h.svc.List(c, search); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessSlice(c, accounts)
	}
	response.SuccessSlice(c, nil)
}

func (h *AccountHandler) create(c *gin.Context) {
	var position types2.AccountCreateReq
	if err := c.ShouldBindJSON(&position); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if account, err := h.svc.Create(c, &position); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessMap(c, account)
	}
}

func (h *AccountHandler) put(c *gin.Context) {
	var account types2.AccountCreateReq
	if err := c.ShouldBindJSON(&account); err != nil {
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
	h.l.Debug(fmt.Sprintf("修改参数: %+v, 修改id: %d", account, id))
	if newPosition, err := h.svc.Put(c, id, &account); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessMap(c, newPosition)
	}

}

func (h *AccountHandler) delete(c *gin.Context) {
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

func (h *AccountHandler) Name() string {
	return fmt.Sprintf("%s.%s", apps.AppName, apps.AppAccount)
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *AccountHandler) Config() {
	h.l = global.L.Named(apps.AppName).Named(apps.AppAccount).Named("handler")
	h.svc = router.GetLogic(h.Name()).(*logic.AccountLogic)
}

func init() {
	router.RegistryGinRouter(accountHandler)
}

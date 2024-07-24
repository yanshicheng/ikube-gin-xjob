package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/logic"
	"github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	types2 "github.com/yanshicheng/ikube-gin-xjob/apps/users/types"
	"github.com/yanshicheng/ikube-gin-xjob/common/response"
	"github.com/yanshicheng/ikube-gin-xjob/common/types"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"github.com/yanshicheng/ikube-gin-xjob/router"
	"go.uber.org/zap"
)

var _ router.GinService = (*PositionHandler)(nil)
var positionHandler = &PositionHandler{}

type PositionHandler struct {
	l   *zap.Logger
	svc *logic.PositionLogic
}

func (h *PositionHandler) PublicRegistry(gin.IRouter) {

}

// AuthRegistry 注册认证接口
func (h *PositionHandler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group(fmt.Sprintf("%s/%s", users.AppName, users.AppPosition))
	{
		group.GET("/", h.list)
		group.POST("/", h.create)
		group.PUT("/:id", h.put)
		group.DELETE("/:id", h.delete)
	}
}

func (h *PositionHandler) list(c *gin.Context) {
	search := types2.PositionListSearchReq{}
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

func (h *PositionHandler) create(c *gin.Context) {
	var position model.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		global.LSys.Error(fmt.Sprintf("参数绑定失败: %s", err.Error()))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.Create(c, &position); err != nil {
		response.FailedStr(c, err.Error())
		return
	}
	response.SuccessMap(c, position)
}

func (h *PositionHandler) put(c *gin.Context) {
	var position model.Position
	if err := c.ShouldBindJSON(&position); err != nil {
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
	h.l.Debug(fmt.Sprintf("修改参数: %+v, 修改id: %d", position, id))
	if newPosition, err := h.svc.Put(c, id, &position); err != nil {
		response.FailedStr(c, err.Error())
		return
	} else {
		response.SuccessMap(c, newPosition)
	}

}

func (h *PositionHandler) delete(c *gin.Context) {
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

func (h *PositionHandler) Name() string {
	return fmt.Sprintf("%s.%s", users.AppName, users.AppPosition)
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *PositionHandler) Config() {
	h.l = global.L.Named(users.AppName).Named(users.AppPosition).Named("handler")
	h.svc = router.GetLogic(h.Name()).(*logic.PositionLogic)
}

func init() {
	router.RegistryGinRouter(positionHandler)
}

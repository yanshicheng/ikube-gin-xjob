package all

import (
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/upms/handler"
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/upms/logic"
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/upms/model"
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/users/handler"
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/users/logic"
	_ "github.com/yanshicheng/ikube-gin-xjob/apps/users/model"
	// 引入自定义验证器
	_ "github.com/yanshicheng/ikube-gin-xjob/common/validator"
)

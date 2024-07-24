package model

import "github.com/yanshicheng/ikube-gin-xjob/global"

func Register(model ...interface{}) {
	global.M = append(global.M, model...)
}

package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/yanshicheng/ikube-gin-xjob/pkg/mysql"
	"github.com/yanshicheng/ikube-gin-xjob/pkg/redis"
	"github.com/yanshicheng/ikube-gin-xjob/pkg/types"
	"go.uber.org/zap"
)

const JwtKey = `zk5Lp8m+0g7lvvLcnbUHPzQFEsRAmvNIn9tXdx0o80U=`

// 静态目录文件夹
const StaticDir = "static"

var (
	IkubeopsTrans ut.Translator
	C             *types.Config = types.NewDefaultConfig()
	L             *zap.Logger
	LSys          *zap.Logger
	DB            *mysql.IkubeGorm
	RDB           *redis.IkubeRedis
	M             []interface{}
)

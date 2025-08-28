package dal

import (
	"tiktok/biz/dal/db"
	"tiktok/biz/mw/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

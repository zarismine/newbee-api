package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

var (
	CacheUserContactPrefix = "cache:contact:"
	CacheUserMessagePrefix = "cache:message:"
	CacheUserCount         = "cache:user:count:"
	CacheCategory          = "cache:category"
)

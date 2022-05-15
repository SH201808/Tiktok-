// Package redis
//
// @author YangHao
//
// @brief 提供redis数据库的一些操作方法
//
// @date 2022-05-15
//
// @version 0.1
//
package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"tiktok/setting"
)

// RedisDB 是初始化以后全局可使用的Redis客户端
var RedisDB *redis.Client

// Init
//
// @author YangHao
//
// @brief 初始化Redis客户端
//
// @params cfg *setting.RedisConfig: 初始化需要使用的配置设置
//
// @return RedisDB *redis.Client: 返回初始化完成的Redis客户端
// 		   err error: 初始化如果产生错误则会返回该错误
//
// @date 2022-5-15
//
// @version 0.1
//
func Init(cfg *setting.RedisConfig) (RedisDB *redis.Client, err error) {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%d/%d", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		// 返 回 错 误
		panic(err)
	}
	opt.PoolSize = cfg.PoolSize
	RedisDB = redis.NewClient(opt)
	return
}

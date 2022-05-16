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
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"tiktok/setting"
	"time"
)

const (
	KeyIsExistedError    = "this key has already in the redis"
	KeyIsNotExistedError = "can not find the key"
	BadTypeOfValueError  = "incr function only receive int64 or float64"
)

// DB 是初始化以后的Redis客户端
var DB *redis.Client

// ctxDefault 是默认的context，如果没有传入自定义的context则使用该context（什么也不做）
var ctxDefault = context.Background()

// Init
//
// @author YangHao
//
// @brief 初始化Redis客户端
//
// @params cfg *setting.RedisConfig: 初始化需要使用的配置设置
//
// @return DB *redis.Client: 返回初始化完成的Redis客户端
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

// Close
//
// @author YangHao
//
// @brief 关闭Redis客户端，发生错误则会打印在日志中
//
// @date 2022-5-16
//
// @version 0.1
//
func Close() {
	defer func(RedisDB *redis.Client) {
		err := RedisDB.Close()
		if err != nil {
			log.Fatal("a fatal error occurred when closing the DB")
			return
		}
	}(DB)
}

// GetValue
//
// @author YangHao
//
// @brief 获取指定key的value
//
// @params ctx context.Context: 上下文，nil时使用默认的context
//		   key string: 目标key
//
// @return value any: 值
//	 	   err error: 错误类型
//
// @date 2022-5-16
//
// @version 0.1
//
func GetValue(ctx context.Context, key string) (value any, err error) {
	if ctx == nil {
		ctx = ctxDefault
	}
	value, err = DB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New(KeyIsNotExistedError)
		}
		return nil, err
	}
	return
}

// SetValue
//
// @author YangHao
//
// @brief 根据需要来set一个key的value进入redis
//
// @params ctx context.Context: 上下文，nil时使用默认的context
//		   key string: 目标key
//		   value any: 值
//         writeOn any: 是否覆盖。当此项为真时，set操作会覆盖已经存在的key，给它设定新的值value
//		   expiredTime time.Duration: 给key设定的过期时间
//
// @return err error: 错误类型
//
// @date 2022-5-16
//
// @version 0.1
//
func SetValue(ctx context.Context, key string, value any, writeOn bool, expiredTime time.Duration) (err error) {
	if writeOn {
		err = DB.SetEX(ctx, key, value, expiredTime).Err()
	} else {
		err = DB.SetNX(ctx, key, value, expiredTime).Err()
		if err == redis.Nil {
			return errors.New(KeyIsExistedError)
		}
	}
	return err
}

// IncrValue
//
// @author YangHao
//
// @brief 目标key的值增加指定数额
//
// @params ctx context.Context: 上下文，nil时使用默认的context
//		   key string: 目标key
//		   value any: 值，只接受int64和float64
//
// @return err error: 错误类型
//
// @date 2022-5-16
//
// @version 0.1
//
func IncrValue(ctx context.Context, key string, value any) (err error) {
	if ctx == nil {
		ctx = ctxDefault
	}

	switch value.(type) {
	case int64:
		if value == 1 {
			err = DB.Incr(ctx, key).Err()
		} else {
			err = DB.IncrBy(ctx, key, value.(int64)).Err()
		}
		return err
	case float64:
		err = DB.IncrByFloat(ctx, key, value.(float64)).Err()
		return err
	}

	return errors.New(BadTypeOfValueError)
}

// DecrValue
//
// @author YangHao
//
// @brief 目标key的值减少指定数额
//
// @params ctx context.Context: 上下文，nil时使用默认的context
//		   key string: 目标key
//		   value int64: 值
//
// @return err error: 错误类型
//
// @date 2022-5-16
//
// @version 0.1
//
func DecrValue(ctx context.Context, key string, value int64) (err error) {
	if ctx == nil {
		ctx = ctxDefault
	}

	if value == 1 {
		err = DB.Decr(ctx, key).Err()
	} else {
		err = DB.DecrBy(ctx, key, value).Err()
	}

	return err
}

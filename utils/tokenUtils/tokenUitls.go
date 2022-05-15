// Package tokenUtils
//
// @author YangHao
//
// @brief 提供Token的生成与解析功能
//
// @date 2022-05-15
//
// @version 0.1
package tokenUtils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"tiktok/setting"
	"time"
)

const (
	// TokenExpiredTime 的单位是以秒计算的
	TokenExpiredTime  = 600
	TokenExpiredError = "token is expired"
	TokenIllegalError = "token is illegal"
	TokenHandleError  = "can not handle this token, maybe it is not a valid token"
)

// TokenClaims 自定义的payload结构
type TokenClaims struct {
	Data any `json:"data"`
	jwt.RegisteredClaims
}

// CreateToken
//
// @author YangHao
//
// @brief 生成一个新的Token，采用了HS256算法
//	其中payload结构采用的是自定义的 TokenClaims
// 	除了需要放置的数据以外只存放了签发时间和生效时间
//
// @params data any: 需要放在payload中的数据
//
// @return token string: 生成的token
// 		   err error: 如果生成token期间发生错误，则返回发生的错误
//
// @date 2022-5-15
//
// @version 0.1
//
func CreateToken(data any) (token string, err error) {
	claim := TokenClaims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	originToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = originToken.SignedString([]byte(setting.Conf.TokenConfig.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseToken
//
// @author YangHao
//
// @brief 对传入的Token进行解析
//
// @params token string: 需要解析的token
//
// @return data any: token中的payload存放的数据
// 		   err error: 返回该token存在的问题
//
// @date 2022-5-15
//
// @version 0.1
//
func ParseToken(token string) (data any, err error) {
	tokenStruct, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.Conf.TokenConfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenStruct.Claims.(*TokenClaims)
	if ok && tokenStruct.Valid {
		if claims.NotBefore.After(time.Now()) {
			return nil, errors.New(TokenIllegalError)
		}
		if time.Now().Unix()-claims.IssuedAt.Unix() > TokenExpiredTime {
			return nil, errors.New(TokenExpiredError)
		}
		return claims.Data, nil
	}
	return nil, errors.New(TokenHandleError)
}

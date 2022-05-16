// Package cryptoUtils
//
// @author YangHao
//
// @brief 提供对数据加密的方法封装，采用AES加密
//
// @date 2022-05-15
//
// @version 0.1
package cryptoUtils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"tiktok/setting"
)

const (
	AesKeyError            = "crypto key is illegal, need to change the config files"
	DecodeStringError      = "there is a problem when decoding the encrypted string"
	UnknownEncodeTypeError = "switch encode_type with an unknown value, it will be encoded as in Standard mode"
	UnknownDecodeTypeError = "switch decode_type with an unknown value"

	UrlOrFile = iota
	Standard
)

// Encrypt
//
// @author YangHao
//
// @brief 对数据进行加密操作，采用AES加密
//
// @params originData []byte: 需要加密的数据的字节数组
// 		   encodeMode int: 加密模式
//
// @return data string: 加密后的数据
// 		   err error: 加密如果出错则返回一个错误，根据错误类型进行处理
//
// @date 2022-5-15
//
// @version 0.1
//
func Encrypt(originData []byte, encodeMode int) (data string, err error) {
	key := []byte(setting.Conf.Key)

	block, e := aes.NewCipher(key)
	if e != nil {
		return "", errors.New(AesKeyError)
	}
	blockSize := block.BlockSize()

	originData = textPKCS7Padding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	encrypted := make([]byte, len(originData))
	blockMode.CryptBlocks(encrypted, originData)

	if encodeMode == Standard {
		return base64.StdEncoding.EncodeToString(encrypted), nil
	}
	if encodeMode == UrlOrFile {
		return base64.URLEncoding.EncodeToString(encrypted), nil
	}

	return base64.StdEncoding.EncodeToString(encrypted), errors.New(UnknownEncodeTypeError)
}

// Decrypt
//
// @author YangHao
//
// @brief 对数据进行解密操作，采用AES解密
//
// @params encryptedData string: 需要解密的字符串
// 		   decodeMode int: 解密模式
//
// @return data []data: 解密后的数据
// 		   err error: 解密如果出错则返回一个错误，根据错误类型进行处理
//
// @date 2022-5-15
//
// @version 0.1
//
func Decrypt(encryptedData string, decodeMode int) (data []byte, err error) {
	key := []byte(setting.Conf.Key)

	var encryptedByte []byte
	var e error

	if decodeMode == Standard {
		encryptedByte, e = base64.StdEncoding.DecodeString(encryptedData)
	} else if decodeMode == UrlOrFile {
		encryptedByte, e = base64.URLEncoding.DecodeString(encryptedData)
	} else {
		return nil, errors.New(UnknownDecodeTypeError)
	}
	if e != nil {
		return nil, errors.New(DecodeStringError)
	}

	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, errors.New(AesKeyError)
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	origin := make([]byte, len(encryptedByte))
	blockMode.CryptBlocks(origin, encryptedByte)
	origin = textPKCS7UnPadding(origin)

	return origin, nil
}

// textPKCS7Padding
//
// @author YangHao
//
// @brief 对数据原文进行填充
//
// @params ciphertext []byte: 需要填充的数据
// 		   blockSize int: 填充体积
//
// @return 补充字节完毕后的数据
//
// @date 2022-5-15
//
// @version 0.1
//
func textPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// textPKCS7UnPadding
//
// @author YangHao
//
// @brief 对解密的数据正文进行移除填充字节
//
// @params originData []byte: 需要移除填充字节的解密数据
//
// @return 移除填充字节完毕后的数据
//
// @date 2022-5-15
//
// @version 0.1
//
func textPKCS7UnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

// Copyright 2017 guangbo. All rights reserved.

// 加解密散列接口，提供aes，des，md5接口
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/jkkkls/hjing/utils/bytes_cache"
)

// DesEncrypt des加密函数，返回加密后的结果长度是8的倍数
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())

	iv := bytes_cache.Get(des.BlockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := bytes_cache.Get(len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// DesDecrypt des解密函数，传入解密内容长度必须是8的倍数
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := bytes_cache.Get(des.BlockSize)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := bytes_cache.Get(len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt aes加密函数，返回加密后的结果长度是16的倍数
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)

	// iv := bytes_cache.Get( aes.BlockSize)
	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := bytes_cache.Get(len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt aes解密函数，传入解密内容长度必须是16的倍数
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// iv := bytes_cache.Get( aes.BlockSize)
	iv := []byte("0000000000000000")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := bytes_cache.Get(len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

// Md5 md5散列函数
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Uint64() uint64 { return rand.Uint64() }

func Int31() int32 { return rand.Int31() }

func Int63() int64 { return int64(Uint64() >> 1) }

func Int31n(n int32) int32 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int31() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := Int31()
	for v > max {
		v = Int31()
	}
	return v % n
}

func Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int63()
	for v > max {
		v = Int63()
	}
	return v % n
}

// Intn 随机一个小于n的整数，[0, n)
func Intn(n int) int {
	if n == 0 {
		panic("invalid argument to int")
	}

	if n <= 1<<31-1 {
		return int(Int31n(int32(n)))
	}
	return int(Int63n(int64(n)))
}

const randomStr = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GetRandomString 生成指定长度的随机字符串
func GetRandomString(n int) string {
	bytes := []byte(randomStr)
	result := []byte{}

	for i := 0; i < n; i++ {
		result = append(result, bytes[Intn(len(bytes))])
	}
	return string(result)
}

// RandomUint32 随机一个指定区间的数，[min, max]
func RandomUint32(min, max uint32) uint32 {
	return uint32(Intn(int(max-min+1))) + min
}

// RandomInt32 随机一个指定区间的数，[min, max]
func RandomInt32(min, max int32) int32 {
	return int32(Intn(int(max-min+1))) + min
}

// ShuffleUint32Array ...
func ShuffleUint32Array(x *[]uint32, n int) {
	for i := 0; i < n; i++ {
		for j := len(*x) - 1; j >= 0; j-- {
			r := Intn(j + 1)
			(*x)[r], (*x)[j] = (*x)[j], (*x)[r]
		}
	}
}

// ShuffleInt32Array ...
func ShuffleInt32Array(x *[]int32, n int) {
	for i := 0; i < n; i++ {
		for j := len(*x) - 1; j >= 0; j-- {
			r := Intn(j + 1)
			(*x)[r], (*x)[j] = (*x)[j], (*x)[r]
		}
	}
}

// RandLessWight32 计算权重，小值优先
// 等待泛型- -#
func RandLessWight32(weights []uint32) int {
	var weights64 []uint64
	for _, v := range weights {
		weights64 = append(weights64, uint64(v))
	}

	return RandLessWight(weights64)
}

// RandLessWight 计算权重，小值优先
func RandLessWight(weights []uint64) int {
	var (
		max, total uint64
	)
	if len(weights) == 0 {
		panic("invalid argument to RandLessWight")
	} else if len(weights) == 1 {
		return 0
	}
	for i := 0; i < len(weights); i++ {
		if weights[i] > max {
			max = weights[i]
		}
	}
	//重新计算权重
	if max < 10 {
		max += 10
	} else {
		max += max / 5
	}
	for i := 0; i < len(weights); i++ {
		weights[i] = max - weights[i]
	}
	for i := 0; i < len(weights); i++ {
		total += weights[i]
	}

	r := Uint64() % total

	for i := 0; i < len(weights); i++ {
		if r < weights[i] {
			return i
		}
		r -= weights[i]
	}

	return 0
}

package storage

import (
	"math/rand"
	"testing"
	"time"
)

func TestQuestion1(t *testing.T) {
	ParallelPrint()
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func TestQuestion2(t *testing.T) {
	// 正常情况
	PrintSubset("abc")
	// 重复情况
	PrintSubset("aaaaa")
	PrintSubset("abcabcabcabc")
	// 空集情况
	PrintSubset("")
	// 长字符串情况
	PrintSubset(RandomStr(1089))
}

func TestQuestion3(t *testing.T) {
	t.Log(GetZeroNums(5))
	t.Log(GetZeroNums(12))
	t.Log(GetZeroNums(0))
}

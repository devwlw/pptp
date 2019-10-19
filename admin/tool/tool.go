package tool

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"

	"github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.NewV4().String()
}

func NonceStr() string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghiklmnopqrstuvwxyz"
	re := make([]byte, 12)
	for i := 0; i < 12; i++ {
		re[i] = str[rand.Intn(len(str))]
	}
	return string(re)
}

func NonceStrN(n int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghiklmnopqrstuvwxyz"
	re := make([]byte, n)
	for i := 0; i < n; i++ {
		re[i] = str[rand.Intn(len(str))]
	}
	return string(re)
}

func NonceNumberN(n int) string {
	str := "0123456789"
	re := make([]byte, n)
	for i := 0; i < n; i++ {
		re[i] = str[rand.Intn(len(str))]
	}
	return string(re)
}

func SplitStringArr(list []string, size int) [][]string {
	arr := make([][]string, 0)
	a := len(list) / size
	b := len(list) % size
	for i := 0; i < a; i++ {
		arr = append(arr, list[i*size:(i+1)*size])
	}
	if b != 0 {
		arr = append(arr, list[a*size:a*size+b])
	}
	return arr
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

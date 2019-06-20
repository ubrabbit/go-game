package common

import (
	"math/rand"
	"time"
)

const (
	G_RANDOM_DIGIT = 0
	G_RANDOM_LOWER = iota
	G_RANDOM_UPPER = iota
	G_RANDOM_ALL   = iota
)

/**
* size 随机码的位数
* kind 0    // 纯数字
       1    // 小写字母
       2    // 大写字母
       3    // 数字、大小写字母
*/
func RandomString(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > G_RANDOM_UPPER || kind < G_RANDOM_DIGIT
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(G_RANDOM_ALL)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}

	return string(result)
}

func RandomList(s []string) string {
	l := len(s)
	if l <= 0 {
		return ""
	}
	idx := rand.Intn(l)
	return s[idx]
}

//  start <= x <= end
func RandomInt(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(end - start + 1)
	return start + value
}

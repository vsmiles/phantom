package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Randomfloat generates a random integer between min and max
func Randomfloat(min, max int64) float64 {
	return float64(min+rand.Int63n(max-min)) + rand.Float64()
}

func randomNotRepeatingInt(start int64, end int64, count int64) []int64 {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int64, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < int(count) {
		//生成随机数
		num := r.Int63n((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder // 字符串生成器
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUser generates a random owner name
func RandomUser() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@vmail.com", RandomString(6))
}

// RandomDate generates a random time
func RandomDate() time.Time {
	return time.Now().AddDate(0, 0, -int(RandomInt(0, 30)))
}

func RandomGenres() []string {
	allGenres := []string{"Short", "Western", "Drama", "Fantasy", "Animation", "Comedy", "Crime"}
	n1 := RandomInt(1, 4)
	n2 := randomNotRepeatingInt(0, 7, n1)
	var genres []string
	for i := 0; i < int(n1); i++ {
		genres = append(genres, allGenres[n2[i]])
	}
	return genres
}

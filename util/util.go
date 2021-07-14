package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string { //随机n位字符串
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n) //创建一个byte数组，长度为n，名为result

	rand.Seed(time.Now().Unix()) //根据时间戳获得随机数种子
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))] //rand是一个库，调用Intn方法，这个方法随机一个数
	}
	return string(result)
}

package main

import "fmt"

// 在Go中，如果我们能够保证从不会向一个通道发送数据，那么有一个简单的方法来判断此通道是否已经关闭。
type T int

func IsClosed(ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func main() {
	c := make(chan T)
	fmt.Println(IsClosed(c)) // false
	close(c)
	fmt.Println(IsClosed(c)) // true
}
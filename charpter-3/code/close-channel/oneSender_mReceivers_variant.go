package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
有时，数据通道（dataCh）的关闭请求需要由某个第三方协程发出。
对于这种情形，我们可以使用一个额外的信号通道来通知唯一的发送者关闭数据通道（dataCh）。
 */
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const Max = 100000
	const NumReceivers = 100
	const NumThirdParties = 15

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int)
	// 传出发起关闭datach消息
	closing := make(chan struct{}) // 信号通道
	// 传输datach已关闭消息
	closed := make(chan struct{})

	// 此stop函数可以被安全地多次调用。
	// 防止第三方在网络请求时出现超时重试，保证幂等。
	// 上述代码中的stop函数中使用的技巧偷自Roger Peppe在此贴中的一个留言。
	stop := func() {
		select {
		case closing<-struct{}{}:
			<-closed // 如果closed已关闭，不会阻塞直接返回
		case <-closed:
		}
	}

	// 一些第三方协程
	for i := 0; i < NumThirdParties; i++ {
		go func() {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop()
		}()
	}

	// 发送者
	go func() {
		defer func() {
			close(closed)
			close(dataCh)
		}()

		for {
			select{
			case <-closing: return
			default:
			}

			select{
			case <-closing: return
			case dataCh <- rand.Intn(Max):
			}
		}
	}()

	// 接收者
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}
package main

import (
	"fmt"
	"sync"
	"time"
)

/**
编写一个程序，使用通道实现两个协程之间的通信。
一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
*/

func method1() {

	ch := make(chan int)

	var wg sync.WaitGroup

	wg.Add(2)
	//生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 0; i <= 10; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}

	}()

	// 消费者协程
	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
		}
		fmt.Println("通道已关闭，接收完毕")
	}()

	// 等待两个协程都完成
	wg.Wait()
	fmt.Println("程序结束")

}

/**
实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
*/

func method2() {

	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	//生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 0; i <= 100; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}

	}()

	//消费者协程
	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
		}
		fmt.Println("通道已关闭，接收完毕")

	}()

	wg.Wait()
	fmt.Println("程序结束")

}

func main() {
	//method1()
	method2()
}

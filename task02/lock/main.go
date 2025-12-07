package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/**
编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值
*/

var x1 int
var wg sync.WaitGroup
var lock sync.Mutex

func add1() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		x1 += 1
		lock.Unlock()
	}
	wg.Done()
}

/*
*
使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/
var x2 int64

func add2() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt64(&x2, 1)
	}
	wg.Done()
}

func main() {
	/*wg.Add(10)
	// 启动10个协程
	for i := 0; i < 10; i++ {
		go add1()
	}
	wg.Wait()
	fmt.Println("最终x1结果:", x1)*/
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go add2()
	}
	wg.Wait()
	fmt.Println("最终x2结果:", x2) // 10000
}

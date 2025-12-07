package main

import (
	"fmt"
	"sync"
	"time"
)

func js() {

	for i := 1; i <= 10; i++ {
		if i%1 == 1 {
			fmt.Println("奇数:", i)

		}
	}
}
func os() {
	for i := 2; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Println("偶数:", i)
		}
	}

}

/*
*
编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
*/
func method1() {
	fmt.Println("=== 启动两个协程 ===")
	jsChan := make(chan bool)
	osChan := make(chan bool)

	// 启动奇数协程
	go func() {
		defer close(osChan) // 奇数协程结束时关闭偶数通道
		for i := 1; i <= 10; i += 2 {
			<-jsChan // 等待信号
			//value := <-jsChan  // 完整写法：接收值并赋值给变量
			//<-jsChan           // 简写：只接收，不保存值（丢弃接收的值）--阻塞的写法
			fmt.Printf("奇数: %d\n", i)
			if i < 9 { // // 最后一个奇数不通知
				osChan <- true // 通知偶数协程
			}
		}
	}()

	// 启动偶数协程
	go func() {
		defer close(jsChan) // 偶数协程结束时关闭奇数通道
		for i := 2; i <= 10; i += 2 {
			<-osChan
			fmt.Printf("偶数: %d\n", i)
			if i < 10 { // 最后一个偶数不通知
				jsChan <- true // 通知奇数协程
			}
		}
	}()

	jsChan <- true // 启动第一个协程（奇数）
	time.Sleep(time.Millisecond * 100)
	fmt.Println()

}

/**
设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
*/

func method2() {
	fmt.Println("=== 任务调度器 ===")
	task1 := func() {
		time.Sleep(time.Second * 1)
		fmt.Println("任务1: 处理用户数据")
	}

	task2 := func() {
		time.Sleep(time.Second * 3)
		fmt.Println("任务2: 发送邮件")
	}

	task3 := func() {
		time.Sleep(time.Second * 3)
		fmt.Println("任务3: 备份数据库")
	}

	//将任务放入切片
	tasks := []func(){task1, task2, task3}

	//创建WaitGroup用于等待所有协程
	var wg sync.WaitGroup
	// 启动协程执行每个任务
	for i, task := range tasks {
		wg.Add(1) // 每次循环增加一个等待任务

		go func(id int, t func()) {
			defer wg.Done() // 任务完成后标记完成

			start := time.Now()
			fmt.Printf("任务%d 开始执行...\n", id+1)

			t() // 执行任务

			duration := time.Since(start)
			fmt.Printf("任务%d 执行完成，耗时: %v\n", id+1, duration)
		}(i, task)

		//等待所有任务完成
		fmt.Println("等待所有任务完成...")
		wg.Wait()
		fmt.Println("所有任务执行完毕！")

	}

}

func main() {

	/*go js()
	go os()
	fmt.Println("main  goroutine done")
	time.Sleep(5 * time.Second)*/

	//method1()
	method2()

}

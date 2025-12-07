package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
*
题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
Gorm 实现
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，执行一个复杂的查询，例如查询价格大于 50 元的书籍
，并将结果映射到 Book 结构体切片中，确保类型安全
*/
type Book struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Title  string  `gorm:"type:varchar(200);not null"`
	Author string  `gorm:"type:varchar(100);not null"`
	Price  float64 `gorm:"type:decimal(10,2);not null"`
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("表结构迁移失败:", err)
	}

	fmt.Println("数据库连接成功，表结构已准备")

	// 创建测试数据
	createSampleData(db)

	// 条件查询
	fmt.Println("=== 示例1：价格大于50元的书籍 ===")
	var expensiveBooks []Book
	err = db.Where("price > ?", 50).Find(&expensiveBooks).Error
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 本价格大于50元的书籍:\n", len(expensiveBooks))
		for _, book := range expensiveBooks {
			fmt.Printf("《%s》- %s (%.2f元)\n", book.Title, book.Author, book.Price)
		}
	}
}

func createSampleData(db *gorm.DB) {
	// 清空表
	db.Exec("DELETE FROM books")

	books := []Book{
		{Title: "活着", Author: "余华", Price: 39.9},
		{Title: "许三观卖血记", Author: "余华", Price: 45.0},
		{Title: "兄弟", Author: "余华", Price: 68.5},
		{Title: "围城", Author: "钱钟书", Price: 55.0},
		{Title: "百年孤独", Author: "加西亚·马尔克斯", Price: 88.0},
		{Title: "小王子", Author: "安托万·德·圣-埃克苏佩里", Price: 29.9},
		{Title: "三体", Author: "刘慈欣", Price: 93.0},
		{Title: "流浪地球", Author: "刘慈欣", Price: 42.5},
		{Title: "白夜行", Author: "东野圭吾", Price: 59.8},
		{Title: "解忧杂货店", Author: "东野圭吾", Price: 49.9},
		{Title: "哈利·波特与魔法石", Author: "J.K.罗琳", Price: 75.0},
		{Title: "1984", Author: "乔治·奥威尔", Price: 36.0},
	}

	result := db.Create(&books)
	if result.Error != nil {
		log.Println("插入示例数据失败:", result.Error)
	} else {
		fmt.Println("已创建示例书籍数据")
	}
}

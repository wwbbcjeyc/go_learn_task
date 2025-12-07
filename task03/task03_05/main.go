package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，
其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。
*/

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Nickname string `gorm:"type:varchar(50)"`
	Avatar   string `gorm:"type:varchar(255)"`
	Bio      string `gorm:"type:text"`
	IsActive bool   `gorm:"default:true"`

	// 一对多关系
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

// Post 文章模型
type Post struct {
	BaseModel
	Title     string `gorm:"type:varchar(200);not null"`
	Content   string `gorm:"type:longtext;not null"`
	Summary   string `gorm:"type:text"`
	Status    string `gorm:"type:varchar(20);default:'draft';index"` // draft, published, archived
	ViewCount int    `gorm:"default:0"`
	LikeCount int    `gorm:"default:0"`

	// 外键
	UserID uint `gorm:"index;not null"`
	// 关系
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content   string `gorm:"type:text;not null"`
	IsPublic  bool   `gorm:"default:true"`
	IPAddress string `gorm:"type:varchar(45)"`

	// 外键
	UserID uint `gorm:"index;not null"`
	PostID uint `gorm:"index;not null"`
	// 关系
	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 显示SQL日志
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 获取底层SQL数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败:", err)
	}
	defer sqlDB.Close()

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	fmt.Println("开始创建数据库表...")
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}

	fmt.Println("数据库表创建成功!")

	// 验证表结构
	showTableInfo(db)

	// 创建测试数据
	createTestData(db)
}

// 显示表结构信息
func showTableInfo(db *gorm.DB) {
	var tableNames []string
	db.Raw("SHOW TABLES").Scan(&tableNames)

	fmt.Println("\n=== 数据库表列表 ===")
	for _, table := range tableNames {
		fmt.Printf("- %s\n", table)

		// 显示表结构
		var columns []struct {
			Field   string `gorm:"column:Field"`
			Type    string `gorm:"column:Type"`
			Null    string `gorm:"column:Null"`
			Key     string `gorm:"column:Key"`
			Default string `gorm:"column:Default"`
			Extra   string `gorm:"column:Extra"`
		}

		db.Raw(fmt.Sprintf("DESCRIBE %s", table)).Scan(&columns)

		for _, col := range columns {
			fmt.Printf("  %-20s %-20s %-5s %-5s %-10s %s\n",
				col.Field, col.Type, col.Null, col.Key, col.Default, col.Extra)
		}
		fmt.Println()
	}
}

// 创建测试数据
func createTestData(db *gorm.DB) {
	fmt.Println("\n=== 创建测试数据 ===")

	// 创建用户
	users := []User{
		{
			Username: "alice",
			Email:    "alice@example.com",
			Password: "$2a$10$xxx", // 假设是加密后的密码
			Nickname: "Alice",
			Bio:      "热爱技术的开发者",
		},
		{
			Username: "bob",
			Email:    "bob@example.com",
			Password: "$2a$10$yyy",
			Nickname: "Bob",
			Bio:      "写作爱好者",
		},
	}

	// 批量创建用户
	result := db.Create(&users)
	if result.Error != nil {
		log.Println("创建用户失败:", result.Error)
	} else {
		fmt.Printf("创建了 %d 个用户\n", result.RowsAffected)
	}

	// 创建文章
	posts := []Post{
		{
			Title:   "Go语言入门教程",
			Content: "Go语言是一种简洁、高效、并发友好的编程语言...",
			Summary: "介绍Go语言的基本特性和入门知识",
			Status:  "published",
			UserID:  users[0].ID,
		},
		{
			Title:   "GORM使用指南",
			Content: "GORM是Go语言中非常流行的ORM框架...",
			Summary: "详细介绍GORM的各种功能和使用方法",
			Status:  "published",
			UserID:  users[0].ID,
		},
		{
			Title:   "我的技术博客心得",
			Content: "分享我写技术博客的经验和心得...",
			Summary: "如何写出高质量的技术博客",
			Status:  "draft",
			UserID:  users[1].ID,
		},
	}

	result = db.Create(&posts)
	if result.Error != nil {
		log.Println("创建文章失败:", result.Error)
	} else {
		fmt.Printf("创建了 %d 篇文章\n", result.RowsAffected)
	}

	// 创建评论
	comments := []Comment{
		{
			Content:   "这篇文章写得太好了，学到了很多！",
			UserID:    users[1].ID,
			PostID:    posts[0].ID,
			IPAddress: "192.168.1.100",
		},
		{
			Content:   "期待作者写更多关于GORM的内容",
			UserID:    users[0].ID,
			PostID:    posts[1].ID,
			IPAddress: "192.168.1.101",
		},
		{
			Content:   "非常有帮助，谢谢分享！",
			UserID:    users[1].ID,
			PostID:    posts[0].ID,
			IPAddress: "192.168.1.102",
		},
	}

	result = db.Create(&comments)
	if result.Error != nil {
		log.Println("创建评论失败:", result.Error)
	} else {
		fmt.Printf("创建了 %d 条评论\n", result.RowsAffected)
	}

	fmt.Println("\n=== 数据关联查询示例 ===")

	// 查询用户及其文章
	var userWithPosts User
	db.Preload("Posts").First(&userWithPosts, users[0].ID)
	fmt.Printf("用户 %s 的文章数量: %d\n",
		userWithPosts.Username, len(userWithPosts.Posts))

	// 查询文章及其评论
	var postWithComments Post
	db.Preload("Comments").Preload("User").First(&postWithComments, posts[0].ID)
	fmt.Printf("文章《%s》的评论数量: %d, 作者: %s\n",
		postWithComments.Title, len(postWithComments.Comments),
		postWithComments.User.Nickname)
}

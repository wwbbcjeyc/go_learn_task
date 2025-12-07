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
题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

type BaseModel struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null"`
	Email    string `gorm:"size:100;uniqueIndex;not null"`
	Nickname string `gorm:"size:50"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	BaseModel
	Title     string    `gorm:"size:200;not null"`
	Content   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"size:20;default:'draft';index"`
	ViewCount int       `gorm:"default:0"`
	UserID    uint      `gorm:"index;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Comments  []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	BaseModel
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"index;not null"`
	PostID  uint   `gorm:"index;not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Post    Post   `gorm:"foreignKey:PostID"`
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 减少日志输出
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 创建测试数据
	createTestData(db)

	fmt.Println("=== 1. 查询用户所有文章及评论 ===")
	queryUserPostsWithCommentsExample(db)

	fmt.Println("\n=== 2. 查询评论最多的文章 ===")
	queryMostCommentedPostExample(db)

	fmt.Println("\n=== 3. 其他关联查询示例 ===")
	otherAssociationQueries(db)
}

func queryUserPostsWithCommentsExample(db *gorm.DB) {
	// 查询用户ID为1的所有文章及评论
	var user User
	err := db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "published") // 只查询已发布的文章
	}).Preload("Posts.Comments").
		Preload("Posts.Comments.User").
		First(&user, 1).Error

	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}

	fmt.Printf("用户 '%s' 的已发布文章:\n", user.Nickname)
	for i, post := range user.Posts {
		fmt.Printf("\n%d. 《%s》 (浏览量: %d)\n", i+1, post.Title, post.ViewCount)
		fmt.Printf("   评论数: %d\n", len(post.Comments))

		if len(post.Comments) > 0 {
			fmt.Println("   评论列表:")
			for j, comment := range post.Comments {
				fmt.Printf("     %d. %s - %s\n",
					j+1, comment.Content, comment.User.Nickname)
			}
		}
	}
}

func queryMostCommentedPostExample(db *gorm.DB) {
	// 使用子查询获取评论最多的文章
	var postWithStats struct {
		Post
		CommentCount int `gorm:"column:comment_count"`
	}

	// 方法1: 使用JOIN和GROUP BY
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Where("posts.status = ?", "published").
		Group("posts.id").
		Order("comment_count DESC").
		Preload("User").
		First(&postWithStats).Error

	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}

	fmt.Printf("评论最多的文章:\n")
	fmt.Printf("标题: 《%s》\n", postWithStats.Title)
	fmt.Printf("作者: %s\n", postWithStats.User.Nickname)
	fmt.Printf("评论数量: %d\n", postWithStats.CommentCount)
	fmt.Printf("发布日期: %s\n",
		postWithStats.CreatedAt.Format("2006-01-02 15:04"))

	// 获取该文章的前5条评论
	var recentComments []Comment
	db.Where("post_id = ?", postWithStats.ID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentComments)

	if len(recentComments) > 0 {
		fmt.Println("最近评论:")
		for i, comment := range recentComments {
			fmt.Printf("  %d. %s\n", i+1, comment.Content)
		}
	}
}

func otherAssociationQueries(db *gorm.DB) {
	fmt.Println("=== 查询每篇文章的评论统计 ===")

	type PostStats struct {
		PostID       uint   `gorm:"column:post_id"`
		Title        string `gorm:"column:title"`
		Author       string `gorm:"column:author"`
		CommentCount int    `gorm:"column:comment_count"`
	}

	var stats []PostStats

	db.Raw(`
        SELECT 
            p.id as post_id,
            p.title,
            u.nickname as author,
            COUNT(c.id) as comment_count
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN comments c ON p.id = c.post_id
        WHERE p.status = 'published'
        GROUP BY p.id
        ORDER BY comment_count DESC, p.created_at DESC
    `).Scan(&stats)

	fmt.Println("文章评论统计:")
	for _, stat := range stats {
		fmt.Printf("《%s》- %s (评论: %d)\n",
			stat.Title, stat.Author, stat.CommentCount)
	}

	fmt.Println("\n=== 查询用户评论最多的文章 ===")
	var userCommentStats struct {
		UserID   uint   `gorm:"column:user_id"`
		Username string `gorm:"column:username"`
		PostID   uint   `gorm:"column:post_id"`
		Title    string `gorm:"column:title"`
		Comments int    `gorm:"column:comment_count"`
	}

	db.Raw(`
        SELECT 
            u.id as user_id,
            u.username,
            p.id as post_id,
            p.title,
            COUNT(c.id) as comment_count
        FROM users u
        JOIN posts p ON u.id = p.user_id
        LEFT JOIN comments c ON p.id = c.post_id
        GROUP BY u.id, p.id
        ORDER BY u.id, comment_count DESC
    `).Scan(&userCommentStats)

	fmt.Printf("用户 %s 评论最多的文章: 《%s》 (%d条评论)\n",
		userCommentStats.Username, userCommentStats.Title, userCommentStats.Comments)
}

func createTestData(db *gorm.DB) {
	// 清空数据
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	// 创建用户
	users := []User{
		{Username: "tech_guru", Email: "guru@tech.com", Nickname: "技术达人"},
		{Username: "writer_john", Email: "john@write.com", Nickname: "作家约翰"},
		{Username: "coder_lisa", Email: "lisa@code.com", Nickname: "程序员丽莎"},
	}
	db.Create(&users)

	// 创建文章
	posts := []Post{
		{Title: "Go语言并发编程", Content: "详细讲解Go语言并发...", Status: "published", UserID: users[0].ID},
		{Title: "GORM高级技巧", Content: "深入探讨GORM的使用...", Status: "published", UserID: users[0].ID},
		{Title: "微服务架构设计", Content: "微服务架构的实践经验...", Status: "published", UserID: users[0].ID},
		{Title: "我的写作心得", Content: "分享我的写作经验...", Status: "published", UserID: users[1].ID},
		{Title: "小说创作指南", Content: "如何创作一部好的小说...", Status: "draft", UserID: users[1].ID},
		{Title: "React Hooks详解", Content: "React Hooks的使用技巧...", Status: "published", UserID: users[2].ID},
	}
	db.Create(&posts)

	// 创建评论
	comments := []Comment{
		// 文章1的评论（最多评论）
		{Content: "非常实用的教程！", UserID: users[1].ID, PostID: posts[0].ID},
		{Content: "学到了很多，谢谢！", UserID: users[2].ID, PostID: posts[0].ID},
		{Content: "期待更多关于channel的内容", UserID: users[1].ID, PostID: posts[0].ID},
		{Content: "示例代码很清晰", UserID: users[2].ID, PostID: posts[0].ID},
		{Content: "讲得很透彻", UserID: users[1].ID, PostID: posts[0].ID},

		// 文章2的评论
		{Content: "GORM用起来很方便", UserID: users[1].ID, PostID: posts[1].ID},
		{Content: "关联查询部分讲得很好", UserID: users[2].ID, PostID: posts[1].ID},

		// 文章4的评论
		{Content: "写作经验很有帮助", UserID: users[0].ID, PostID: posts[3].ID},
	}
	db.Create(&comments)

	fmt.Println("测试数据创建完成")
}

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
题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，
如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"size:50;uniqueIndex;not null"`
	Nickname  string    `gorm:"size:50"`
	PostCount int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Posts     []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Post struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Title        string    `gorm:"size:200;not null"`
	Content      string    `gorm:"type:text;not null"`
	Status       string    `gorm:"size:20;default:'draft';index"`
	ViewCount    int       `gorm:"default:0"`
	CommentCount int       `gorm:"default:0"`
	HasComments  bool      `gorm:"default:false"`
	UserID       uint      `gorm:"index;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
}

// AfterCreate 钩子 - 创建文章后更新用户文章计数
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	log.Printf("✅ Post AfterCreate: 文章《%s》创建成功，更新用户(ID:%d)文章计数\n",
		p.Title, p.UserID)

	// 使用事务确保一致性
	err = tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).Error

	if err != nil {
		log.Printf("❌ 更新用户文章计数失败: %v", err)
		return err
	}

	log.Printf("✅ 用户(ID:%d)的文章数量已更新\n", p.UserID)
	return nil
}

// BeforeDelete 钩子 - 删除文章前处理
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	log.Printf("⚠️  Post BeforeDelete: 准备删除文章(ID:%d)《%s》\n", p.ID, p.Title)
	return nil
}

// AfterDelete 钩子 - 删除文章后更新用户文章计数
func (p *Post) AfterDelete(tx *gorm.DB) (err error) {
	log.Printf("✅ Post AfterDelete: 文章(ID:%d)删除成功，更新用户文章计数\n", p.ID)

	// 减少用户的文章数量
	err = tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("GREATEST(post_count - 1, 0)")).Error

	if err != nil {
		log.Printf("❌ 更新用户文章计数失败: %v", err)
		return err
	}

	log.Printf("✅ 用户(ID:%d)的文章数量已减少\n", p.UserID)
	return nil
}

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"index;not null"`
	PostID    uint      `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}

// AfterCreate 钩子 - 创建评论后更新文章评论计数
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	log.Printf("✅ Comment AfterCreate: 评论创建成功，更新文章(ID:%d)评论计数\n", c.PostID)

	err = tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Updates(map[string]interface{}{
			"comment_count": gorm.Expr("comment_count + 1"),
			"has_comments":  true,
		}).Error

	if err != nil {
		log.Printf("❌ 更新文章评论计数失败: %v", err)
		return err
	}

	log.Printf("✅ 文章(ID:%d)的评论数量+1，已标记为有评论\n", c.PostID)
	return nil
}

// BeforeDelete 钩子 - 删除评论前检查
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	log.Printf("⚠️  Comment BeforeDelete: 准备删除评论(ID:%d)\n", c.ID)
	return nil
}

// AfterDelete 钩子 - 删除评论后检查并更新文章状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	log.Printf("✅ Comment AfterDelete: 评论(ID:%d)删除成功，检查文章评论状态\n", c.ID)

	// 查询文章的剩余评论数量
	var remainingCount int64
	err = tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&remainingCount).Error

	if err != nil {
		log.Printf("❌ 查询文章评论数量失败: %v", err)
		return err
	}

	log.Printf("📊 文章(ID:%d)剩余评论数量: %d\n", c.PostID, remainingCount)

	// 根据剩余评论数量更新文章状态
	updates := map[string]interface{}{
		"comment_count": remainingCount,
	}

	if remainingCount == 0 {
		updates["has_comments"] = false
		log.Printf("🔄 文章(ID:%d)已更新为无评论状态\n", c.PostID)
	}

	err = tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Updates(updates).Error

	if err != nil {
		log.Printf("❌ 更新文章评论状态失败: %v", err)
		return err
	}

	log.Printf("✅ 文章(ID:%d)评论状态已更新\n", c.PostID)
	return nil
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}

	fmt.Println("🚀 开始演示GORM钩子函数...")

	// 创建测试数据
	createTestData(db)

	// 测试Post钩子
	testPostHooks(db)

	// 测试Comment钩子
	testCommentHooks(db)

	// 显示最终数据状态
	showFinalStatus(db)
}

func createTestData(db *gorm.DB) {
	fmt.Println("\n=== 创建测试数据 ===")

	// 清空数据
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	// 创建用户
	users := []User{
		{Username: "author_li", Nickname: "李作者"},
		{Username: "writer_wang", Nickname: "王作家"},
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatal("创建用户失败:", err)
	}

	fmt.Printf("✅ 创建了 %d 个用户\n", len(users))

	// 验证用户初始状态
	for _, user := range users {
		var dbUser User
		db.First(&dbUser, user.ID)
		fmt.Printf("用户 %s: 文章数量 = %d\n", dbUser.Nickname, dbUser.PostCount)
	}
}

func testPostHooks(db *gorm.DB) {
	fmt.Println("\n=== 测试Post模型钩子 ===")

	// 获取第一个用户
	var user User
	db.First(&user)

	fmt.Printf("用户 %s 的初始文章数量: %d\n", user.Nickname, user.PostCount)

	// 创建文章 - 会触发 AfterCreate 钩子
	fmt.Println("\n1. 创建第一篇文章...")
	post1 := Post{
		Title:   "Go语言学习指南",
		Content: "这是一篇关于Go语言的教程...",
		Status:  "published",
		UserID:  user.ID,
	}

	if err := db.Create(&post1).Error; err != nil {
		log.Printf("创建文章失败: %v", err)
	}

	// 验证用户文章数量已更新
	db.First(&user, user.ID)
	fmt.Printf("创建文章后，用户 %s 的文章数量: %d\n", user.Nickname, user.PostCount)

	// 创建第二篇文章
	fmt.Println("\n2. 创建第二篇文章...")
	post2 := Post{
		Title:   "GORM深入浅出",
		Content: "详细介绍GORM的使用方法...",
		Status:  "published",
		UserID:  user.ID,
	}

	if err := db.Create(&post2).Error; err != nil {
		log.Printf("创建文章失败: %v", err)
	}

	db.First(&user, user.ID)
	fmt.Printf("创建第二篇文章后，用户 %s 的文章数量: %d\n", user.Nickname, user.PostCount)

	// 删除文章 - 会触发 AfterDelete 钩子
	fmt.Println("\n3. 删除一篇文章...")
	if err := db.Delete(&post1).Error; err != nil {
		log.Printf("删除文章失败: %v", err)
	}

	db.First(&user, user.ID)
	fmt.Printf("删除文章后，用户 %s 的文章数量: %d\n", user.Nickname, user.PostCount)
}

func testCommentHooks(db *gorm.DB) {
	fmt.Println("\n=== 测试Comment模型钩子 ===")

	// 获取文章
	var post Post
	db.Where("title = ?", "GORM深入浅出").First(&post)

	fmt.Printf("文章《%s》初始状态:\n", post.Title)
	fmt.Printf("  评论数量: %d\n", post.CommentCount)
	fmt.Printf("  是否有评论: %v\n", post.HasComments)

	// 创建评论 - 会触发 AfterCreate 钩子
	var user User
	db.First(&user)

	fmt.Println("\n1. 创建第一条评论...")
	comment1 := Comment{
		Content: "这篇文章写得太好了！",
		UserID:  user.ID,
		PostID:  post.ID,
	}

	if err := db.Create(&comment1).Error; err != nil {
		log.Printf("创建评论失败: %v", err)
	}

	// 验证文章状态
	db.First(&post, post.ID)
	fmt.Printf("创建评论后，文章《%s》状态:\n", post.Title)
	fmt.Printf("  评论数量: %d\n", post.CommentCount)
	fmt.Printf("  是否有评论: %v\n", post.HasComments)

	// 创建第二条评论
	fmt.Println("\n2. 创建第二条评论...")
	comment2 := Comment{
		Content: "非常详细的教程，感谢分享！",
		UserID:  user.ID,
		PostID:  post.ID,
	}

	if err := db.Create(&comment2).Error; err != nil {
		log.Printf("创建评论失败: %v", err)
	}

	db.First(&post, post.ID)
	fmt.Printf("创建第二条评论后，文章《%s》状态:\n", post.Title)
	fmt.Printf("  评论数量: %d\n", post.CommentCount)
	fmt.Printf("  是否有评论: %v\n", post.HasComments)

	// 删除一条评论 - 会触发 AfterDelete 钩子
	fmt.Println("\n3. 删除一条评论...")
	if err := db.Delete(&comment1).Error; err != nil {
		log.Printf("删除评论失败: %v", err)
	}

	db.First(&post, post.ID)
	fmt.Printf("删除一条评论后，文章《%s》状态:\n", post.Title)
	fmt.Printf("  评论数量: %d\n", post.CommentCount)
	fmt.Printf("  是否有评论: %v\n", post.HasComments)

	// 删除最后一条评论 - 应该将文章标记为无评论
	fmt.Println("\n4. 删除最后一条评论...")
	if err := db.Delete(&comment2).Error; err != nil {
		log.Printf("删除评论失败: %v", err)
	}

	db.First(&post, post.ID)
	fmt.Printf("删除所有评论后，文章《%s》状态:\n", post.Title)
	fmt.Printf("  评论数量: %d\n", post.CommentCount)
	fmt.Printf("  是否有评论: %v\n", post.HasComments)
}

func showFinalStatus(db *gorm.DB) {
	fmt.Println("\n=== 最终数据状态 ===")

	// 显示所有用户
	var users []User
	db.Find(&users)

	fmt.Println("\n用户文章统计:")
	for _, user := range users {
		fmt.Printf("  %s: %d篇文章\n", user.Nickname, user.PostCount)
	}

	// 显示所有文章
	var posts []Post
	db.Preload("User").Find(&posts)

	fmt.Println("\n文章评论统计:")
	for _, post := range posts {
		fmt.Printf("  《%s》(作者: %s): %d条评论, 有评论: %v\n",
			post.Title, post.User.Nickname, post.CommentCount, post.HasComments)
	}
}

// 额外的钩子函数示例
func additionalHookExamples() {
	fmt.Println("\n=== 其他钩子函数示例 ===")

	// Post 的其他钩子
	fmt.Println("1. BeforeSave - 保存前验证")
	fmt.Println("2. AfterSave - 保存后处理")
	fmt.Println("3. BeforeUpdate - 更新前处理")
	fmt.Println("4. AfterFind - 查询后处理")
}

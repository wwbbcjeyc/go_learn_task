package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
*
假设有一个名为 students 的表，包含字段 id （主键，自增）、
name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、
grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/
type Student struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"column:name;type:varchar(100);not null"`
	Age   int    `gorm:"column:age;not null"`
	Grade string `gorm:"column:grade;not null"`
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	db := ConnectDB()
	err := db.AutoMigrate(&Student{})
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectDB 连接数据库
func ConnectDB() *gorm.DB {
	//配置MySQL连接参数
	username := "root"   //账号
	password := "123456" //密码
	host := "127.0.0.1"  //数据库地址，可以是Ip或者域名
	port := 3306         //数据库端口
	Dbname := "test"     //数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	return db
}

// 1. 插入新记录：学生姓名为 "张三"，年龄为 20，年级为 "三年级"
func createStudent(db *gorm.DB) {
	student := Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	}
	result := db.Create(&student)
	if result.Error != nil {
		fmt.Printf("插入失败: %v\n", result.Error)
		return
	}
	fmt.Printf("插入成功，ID: %d\n", student.ID)
}

// 2. 查询所有年龄大于 18 岁的学生信息
func queryStudents(db *gorm.DB) {
	var students []Student

	result := db.Where("age > ?", 18).First(&students)
	if result.Error != nil {
		fmt.Printf("查询失败::%v\n", result.Error)
	}
	fmt.Printf("找到 %d 条记录:\n", result.RowsAffected)

	for _, student := range students {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n",
			student.ID, student.Name, student.Age, student.Grade)
	}

}

// 3. 将姓名为 "张三" 的学生年级更新为 "四年级"
func updateStudent(db *gorm.DB) {
	result := db.Model(&Student{}).
		Where("name = ?", "张三").
		Update("grade", "四年级")

	if result.Error != nil {
		fmt.Printf("更新失败: %v\n", result.Error)
		return
	}

	fmt.Printf("更新成功，影响 %d 条记录\n", result.RowsAffected)
}

// 4. 删除年龄小于 15 岁的学生记录
func deleteStudents(db *gorm.DB) {
	result := db.Where("age < ?", 15).Delete(&Student{})

	if result.Error != nil {
		fmt.Printf("删除失败: %v\n", result.Error)
		return
	}

	fmt.Printf("删除成功，影响 %d 条记录\n", result.RowsAffected)
}

func main() {
	db := InitDB()
	// 1. 插入数据
	fmt.Println("=== 1. 插入新记录 ===")
	createStudent(db)

	// 2. 查询数据
	fmt.Println("\n=== 2. 查询年龄大于18的学生 ===")
	queryStudents(db)

	// 3. 更新数据
	fmt.Println("\n=== 3. 更新张三的年级 ===")
	updateStudent(db)

	// 再次查询验证更新
	fmt.Println("\n=== 验证更新结果 ===")
	queryStudents(db)

	// 4. 删除数据
	fmt.Println("\n=== 4. 删除年龄小于15的学生 ===")
	deleteStudents(db)

}

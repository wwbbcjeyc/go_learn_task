package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Employee struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"type:varchar(100);not null"`
	Department string `gorm:"type:varchar(100);not null"`
	Salary     int    `gorm:"not null"`
}

/*
*
有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码,查询 employees 表中所有部门为 "技术部" 的员工信息，
并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，查询 employees 表中工资最高的员工信息，
并将结果映射到一个 Employee 结构体中。
*/
func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	db.AutoMigrate(&Employee{})

	//插入示例数据用于测试
	createSampleData(db)

	// 1. 查询部门为"技术部"的员工
	fmt.Println("=== 1. 查询技术部员工 ===")
	techEmployees, err := getEmployeesByDepartment(db, "技术部")
	if err != nil {
		log.Println("查询失败:", err)
	}

	// 2. 查询工资最高的员工
	fmt.Println("\n=== 2. 查询工资最高的员工 ===")
	highestPaidEmployee, err := getHighestPaidEmployee(db)
	if err != nil {
		log.Println("查询失败:", err)
	}

	// 避免未使用变量警告
	_ = techEmployees
	_ = highestPaidEmployee
}

func getEmployeesByDepartment(db *gorm.DB, department string) ([]Employee, error) {
	var employees []Employee

	// 使用 Where 进行条件查询
	result := db.Where("department = ?", department).Find(&employees)

	if result.Error != nil {
		return nil, fmt.Errorf("查询部门员工失败: %v", result.Error)
	}

	fmt.Printf("找到 %d 名%s员工:\n", result.RowsAffected, department)
	for i, emp := range employees {
		fmt.Printf("%d. ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
			i+1, emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	return employees, nil
}

func getHighestPaidEmployee(db *gorm.DB) (*Employee, error) {
	var employee Employee

	// 按工资降序排列，取第一条记录
	result := db.Order("salary DESC").First(&employee)

	if result.Error != nil {
		return nil, fmt.Errorf("查询最高工资员工失败: %v", result.Error)
	}

	fmt.Printf("工资最高的员工: ID: %d, 姓名: %s, 部门: %s, 工资: %d\n",
		employee.ID, employee.Name, employee.Department, employee.Salary)

	return &employee, nil
}

// 创建示例数据的函数
func createSampleData(db *gorm.DB) {
	// 清空表
	db.Exec("DELETE FROM employees")

	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 8000},
		{Name: "李四", Department: "技术部", Salary: 9500},
		{Name: "王五", Department: "销售部", Salary: 7000},
		{Name: "赵六", Department: "技术部", Salary: 12000}, // 最高工资
		{Name: "孙七", Department: "人事部", Salary: 6000},
		{Name: "周八", Department: "技术部", Salary: 8500},
	}

	// 批量插入
	result := db.Create(&employees)
	if result.Error != nil {
		log.Println("插入示例数据失败:", result.Error)
	} else {
		fmt.Println("已创建示例数据")
	}
}

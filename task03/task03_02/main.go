package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和
transactions 表（包含字段 id 主键，
from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，
需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
使用gorm
*/

type Account struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	AccountNo string    `gorm:"type:varchar(50);uniqueIndex;not null"` // 账户号
	Balance   float64   `gorm:"type:decimal(15,2);not null;default:0.00"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	FromAccountID uint      `gorm:"not null;index"`
	ToAccountID   uint      `gorm:"not null;index"`
	Amount        float64   `gorm:"type:decimal(15,2);not null"`
	Status        string    `gorm:"type:varchar(20);not null;default:'pending'"`
	TransactionNo string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Remark        string    `gorm:"type:varchar(255)"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	db.AutoMigrate(&Account{}, &Transaction{})

	// 创建测试账户
	createTestAccounts(db)

	// 显示初始余额
	showAccountBalances(db)

	fmt.Println("\n=== 测试正常转账 ===")
	// 测试正常转账
	err = TransferMoneyV2(db, 1, 2, 100.00, "工资转账")
	if err != nil {
		log.Printf("转账失败: %v", err)
	}

	showAccountBalances(db)

	fmt.Println("\n=== 测试余额不足 ===")
	// 测试余额不足
	err = TransferMoneyV2(db, 1, 2, 10000.00, "大额转账")
	if err != nil {
		log.Printf("预期失败: %v", err)
	}

	fmt.Println("\n=== 测试并发转账 ===")
	// 测试并发转账
	testConcurrentTransfers(db)

	// 显示交易记录
	showTransactionHistory(db)
}

// TransferMoneyV2 更健壮的转账实现
func TransferMoneyV2(db *gorm.DB, fromAccountID, toAccountID uint, amount float64, remark string) error {
	// 验证参数
	if err := validateTransferParams(fromAccountID, toAccountID, amount); err != nil {
		return err
	}

	transactionNo := generateTransactionNo()

	return db.Transaction(func(tx *gorm.DB) error {
		// 创建初始交易记录（pending状态）
		transaction := Transaction{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
			Status:        "processing",
			TransactionNo: transactionNo,
			Remark:        remark,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("创建交易记录失败: %w", err)
		}

		// 使用 Raw SQL 进行原子更新，避免并发问题
		sql := `
            UPDATE accounts 
            SET balance = CASE 
                WHEN id = ? AND balance >= ? THEN balance - ?
                WHEN id = ? THEN balance + ?
                ELSE balance
            END
            WHERE id IN (?, ?)
        `

		result := tx.Exec(sql,
			fromAccountID, amount, amount,
			toAccountID, amount,
			fromAccountID, toAccountID)

		if result.Error != nil {
			// 更新交易状态为失败
			tx.Model(&transaction).Update("status", "failed")
			return fmt.Errorf("更新账户余额失败: %w", result.Error)
		}

		// 检查是否成功扣款
		rowsAffected := result.RowsAffected
		if rowsAffected != 2 {
			// 可能扣款失败
			tx.Model(&transaction).Update("status", "failed")
			return fmt.Errorf("转账失败，请检查账户余额")
		}

		// 验证转账结果
		var fromBalance, toBalance float64
		tx.Model(&Account{}).Where("id = ?", fromAccountID).Pluck("balance", &fromBalance)
		tx.Model(&Account{}).Where("id = ?", toAccountID).Pluck("balance", &toBalance)

		// 更新交易状态为成功
		tx.Model(&transaction).Update("status", "success")

		log.Printf("✅ 转账成功: %s, 金额: %.2f, 转出账户余额: %.2f, 转入账户余额: %.2f",
			transactionNo, amount, fromBalance, toBalance)

		return nil
	})
}

// 生成交易流水号
func generateTransactionNo() string {
	return fmt.Sprintf("TXN%013d", time.Now().UnixNano()/1000)
}

// 验证转账参数
func validateTransferParams(fromAccountID, toAccountID uint, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("转账金额必须大于0")
	}

	if fromAccountID == toAccountID {
		return fmt.Errorf("不能向自己转账")
	}

	if amount > 1000000 { // 设置单笔转账限额
		return fmt.Errorf("单笔转账金额不能超过100万")
	}

	return nil
}

// 创建测试账户
func createTestAccounts(db *gorm.DB) {
	// 清空数据
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM accounts")

	accounts := []Account{
		{AccountNo: "622848001001", Balance: 1000.00},
		{AccountNo: "622848001002", Balance: 500.00},
		{AccountNo: "622848001003", Balance: 2000.00},
	}

	if err := db.Create(&accounts).Error; err != nil {
		log.Fatal("创建测试账户失败:", err)
	}

	fmt.Println("✅ 测试账户创建成功")
}

// 显示账户余额
func showAccountBalances(db *gorm.DB) {
	var accounts []Account
	db.Find(&accounts)

	fmt.Println("\n=== 账户余额 ===")
	for _, acc := range accounts {
		fmt.Printf("账户 %d (%s): %.2f 元\n", acc.ID, acc.AccountNo, acc.Balance)
	}
}

// 显示交易历史
func showTransactionHistory(db *gorm.DB) {
	var transactions []Transaction
	db.Order("created_at DESC").Limit(10).Find(&transactions)

	if len(transactions) == 0 {
		fmt.Println("\n暂无交易记录")
		return
	}

	fmt.Println("\n=== 最近交易记录 ===")
	for _, txn := range transactions {
		statusIcon := "✅"
		if txn.Status != "success" {
			statusIcon = "❌"
		}

		fmt.Printf("%s %s: %d → %d, 金额: %.2f, 状态: %s, 时间: %s\n",
			statusIcon, txn.TransactionNo, txn.FromAccountID, txn.ToAccountID,
			txn.Amount, txn.Status, txn.CreatedAt.Format("15:04:05"))
	}
}

// 测试并发转账
func testConcurrentTransfers(db *gorm.DB) {
	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	// 模拟10个并发转账
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 从账户1向账户2转账50元
			err := TransferMoneyV2(db, 1, 2, 50.00, fmt.Sprintf("并发转账测试 %d", index))

			mu.Lock()
			if err != nil {
				failCount++
			} else {
				successCount++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Printf("并发转账结果: 成功 %d 笔, 失败 %d 笔\n", successCount, failCount)
	showAccountBalances(db)
}

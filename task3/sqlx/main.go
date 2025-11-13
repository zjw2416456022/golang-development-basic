/*
	[Sqlx入门]
		题目1：使用SQL扩展库进行查询
			假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
		要求 ：
			编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
			编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
		题目2：实现类型安全映射
			假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
		要求 ：
			定义一个 Book 结构体，包含与 books 表对应的字段。
			编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

func connectDB() (*sqlx.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接初始化失败：%w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败：%w", err)
	}
	// 配置连接池（可选，优化性能）
	db.SetMaxOpenConns(20)   // 最大活跃连接数
	db.SetMaxIdleConns(10)   // 最大空闲连接数
	db.SetConnMaxLifetime(0) // 连接最大生命周期（0 表示无限制）
	log.Println("数据库连接成功")
	return db, nil
}

// 题目1：employees 表查询
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`       // 员工姓名
	Department string  `db:"department"` // 部门
	Salary     float64 `db:"salary"`     // 工资
}

// 查询所有部门为"技术部"的员工，映射到 Employee 切片
func queryTechDepartmentEmployees(db *sqlx.DB) ([]Employee, error) {
	sqlStr := "SELECT id, name, department, salary FROM employees WHERE department = ?"
	// 定义接收结果的切片（强类型，保障类型安全）
	var techEmployees []Employee
	// 参数：切片指针、SQL 语句、查询参数（"技术部" 对应 WHERE 条件）
	err := db.Select(&techEmployees, sqlStr, "技术部")
	if err != nil {
		return nil, fmt.Errorf("查询失败：%w", err)
	}
	return techEmployees, nil
}

// 查询工资最高的员工，映射到单个 Employee 结构体
func queryHighestSalaryEmployee(db *sqlx.DB) (Employee, error) {
	// SQL 语句：按工资降序排序，取第一条（LIMIT 1）
	sqlStr := "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1"
	// 定义接收结果的结构体
	var topEmployee Employee
	err := db.QueryRowx(sqlStr).StructScan(&topEmployee)
	if err != nil {
		// 处理[无数据]场景（ErrNoRows 是合法场景，需单独判断）
		if err == sql.ErrNoRows {
			return Employee{}, fmt.Errorf("无员工数据：%w", err)
		}
		return Employee{}, fmt.Errorf("查询失败：%w", err)
	}
	return topEmployee, nil
}

// 题目2：books 表类型安全查询
type Book struct {
	ID     int             `db:"id"`
	Title  string          `db:"title"`  // 书名
	Author string          `db:"author"` // 作者
	Price  decimal.Decimal `db:"price"`  // 价格
}

// 查询价格大于 50 元的书籍，映射到 Book 切片（类型安全）
func queryBooksPriceGreaterThan50(db *sqlx.DB) ([]Book, error) {
	sqlStr := "SELECT id, title, author, price FROM books WHERE price > ?"
	var expensiveBooks []Book
	err := db.Select(&expensiveBooks, sqlStr, 50.0)
	if err != nil {
		return nil, fmt.Errorf("查询失败：%w", err)
	}
	return expensiveBooks, nil
}

func main() {
	// 1. 建立数据库连接
	db, err := connectDB()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	defer db.Close() // 程序退出前关闭连接，避免资源泄漏

	// 2. 执行题目1查询
	fmt.Println("=== 题目1：查询技术部员工 ===")
	techEmployees, err := queryTechDepartmentEmployees(db)
	if err != nil {
		log.Printf("技术部员工查询失败：%v", err)
	} else {
		for idx, emp := range techEmployees {
			fmt.Printf("第%d个员工：ID=%d, 姓名=%s, 部门=%s, 工资=%.2f\n",
				idx+1, emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}

	fmt.Println("\n=== 题目1：查询工资最高员工 ===")
	topEmp, err := queryHighestSalaryEmployee(db)
	if err != nil {
		log.Printf("工资最高员工查询失败：%v", err)
	} else {
		fmt.Printf("工资最高员工：ID=%d, 姓名=%s, 部门=%s, 工资=%.2f\n",
			topEmp.ID, topEmp.Name, topEmp.Department, topEmp.Salary)
	}

	// 3. 执行题目2查询
	fmt.Println("\n=== 题目2：查询价格>50元的书籍 ===")
	expensiveBooks, err := queryBooksPriceGreaterThan50(db)
	if err != nil {
		log.Printf("高价书籍查询失败：%v", err)
	} else {
		for idx, book := range expensiveBooks {
			fmt.Printf("第%d本书：ID=%d, 书名=%s, 作者=%s, 价格=%.2f\n",
				idx+1, book.ID, book.Title, book.Author, book.Price)
		}
	}
}

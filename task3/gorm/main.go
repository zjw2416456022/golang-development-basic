/*
	[进阶gorm]
		题目1：模型定义
			假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
		要求 ：
			使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
			编写Go代码，使用Gorm创建这些模型对应的数据库表。
		题目2：关联查询
			基于上述博客系统的模型定义。
		要求 ：
			编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
			编写Go代码，使用Gorm查询评论数量最多的文章信息。
		题目3：钩子函数
			继续使用博客系统的模型。
		要求 ：
			为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
			为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
*/

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 题目1：模型定义与表创建
// User 模型（用户）：与 Post 是一对多关系
type User struct {
	gorm.Model          // 内置字段：ID(uint)、CreatedAt、UpdatedAt、DeletedAt
	Username     string `gorm:"size:30;not null;uniqueIndex"`                  // 用户名（唯一）
	Email        string `gorm:"size:100;not null;uniqueIndex"`                 // 邮箱（唯一）
	ArticleCount int    `gorm:"not null;default:0"`                            // 文章数量统计
	Posts        []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // 一对多：一个用户有多篇文章（级联删除）
}

// Post 模型（文章）：与 User 一对多，与 Comment 一对多
type Post struct {
	gorm.Model              // 内置字段
	Title         string    `gorm:"size:100;not null"`              // 文章标题
	Content       string    `gorm:"type:text;not null"`             // 文章内容
	UserID        uint      `gorm:"not null"`                       // 外键：关联 User.ID
	User          User      `gorm:"foreignKey:UserID"`              // 关联 User 模型（反向查询）
	Comments      []Comment `gorm:"foreignKey:PostID"`              // 一对多：一篇文章有多条评论
	CommentStatus string    `gorm:"size:20;not null;default:'有评论'"` // 评论状态：有评论/无评论
	CommentCount  int       `gorm:"-"`                              // 评论数量（仅用于查询结果，不映射数据库字段）
}

// Comment 模型（评论）：与 Post 是一对多关系
type Comment struct {
	gorm.Model        // 内置字段
	Content    string `gorm:"size:500;not null"` // 评论内容
	PostID     uint   `gorm:"not null"`          // 外键：关联 Post.ID
	Post       Post   `gorm:"foreignKey:PostID"` // 关联 Post 模型（反向查询）
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 非必要事务场景关闭默认事务，提升性能
	})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败：%w", err)
	}
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库 ping 失败：%w", err)
	}
	log.Println("数据库连接成功")
	return db, nil
}

// 创建数据库表
func createTables(db *gorm.DB) error {
	// 自动迁移：创建/更新表结构（不会删除已有数据）
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return fmt.Errorf("表创建失败：%w", err)
	}
	log.Println("表创建/更新成功")
	return nil
}

// 题目2：关联查询
// 查询某个用户发布的所有文章及其对应的评论（嵌套预加载）
func getUserArticlesWithComments(db *gorm.DB, userID uint) (*User, error) {
	var user User
	// Preload：预加载关联数据，避免 N+1 查询问题
	// 嵌套 Preload("Posts.Comments")：加载用户的文章，同时加载每篇文章的评论
	err := db.Preload("Posts").Preload("Posts.Comments").First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在：%w", err)
		}
		return nil, fmt.Errorf("查询失败：%w", err)
	}
	return &user, nil
}

// 查询评论数量最多的文章信息
func getMostCommentedPost(db *gorm.DB) (*Post, error) {
	var post Post
	// 逻辑：关联 Comment 表 → 按 post_id 分组 → 统计评论数 → 按评论数降序 → 取第一条
	err := db.Model(&Post{}).
		Select("posts.*, count(comments.id) as comment_count").     // 选择文章所有字段 + 评论数统计
		Joins("left join comments on posts.id = comments.post_id"). // 左连接评论表（避免遗漏无评论的文章）
		Group("posts.id").                                          // 按文章 ID 分组（符合 SQL 标准，避免 group by 报错）
		Order("comment_count desc").                                // 按评论数降序
		Limit(1).                                                   // 取第一条（评论最多）
		Scan(&post).Error                                           // 扫描结果到 Post 结构体（CommentCount 会自动赋值）

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("无文章数据：%w", err)
		}
		return nil, fmt.Errorf("查询失败：%w", err)
	}
	return &post, nil
}

// 题目3：钩子函数

// Post 模型钩子：创建文章时自动更新用户的文章数量（BeforeCreate：创建前执行）
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 逻辑：根据 Post 的 UserID 查询用户 → 文章数量 +1 → 保存用户
	var user User
	// 用 tx 操作保证事务一致性（创建文章和更新用户在同一事务）
	if err := tx.First(&user, p.UserID).Error; err != nil {
		return fmt.Errorf("查询用户失败：%w", err)
	}
	// UpdateColumn：仅更新字段，不触发用户模型的钩子函数（避免循环调用）
	return tx.Model(&user).UpdateColumn("article_count", gorm.Expr("article_count + 1")).Error
}

// Comment 模型钩子：删除评论后检查文章评论数量，为 0 则更新状态（AfterDelete：删除后执行）
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	// 统计当前文章剩余的评论数
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return fmt.Errorf("统计评论数失败：%w", err)
	}
	// 若评论数为 0，更新文章状态为 "无评论"
	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论").Error
	}
	return nil
}

func main() {
	// 1. 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatalf("连接失败：%v", err)
	}

	// 2. 创建表（题目1）
	if err := createTables(db); err != nil {
		log.Fatalf("表创建失败：%v", err)
	}

	// 3. 演示：创建测试数据（用于验证后续功能）
	// 3.1 创建用户
	user := User{Username: "张三", Email: "zhangsan@xxx.com"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("创建用户失败：%v", err)
	}
	log.Printf("创建用户：ID=%d, 用户名=%s", user.ID, user.Username)

	// 创建文章（触发 Post 的 BeforeCreate 钩子，自动更新用户文章数量）
	post1 := Post{Title: "测试文章1", Content: "测试文章1......", UserID: user.ID}
	post2 := Post{Title: "测试文章2", Content: "测试文章2......", UserID: user.ID}
	if err := db.Create(&[]Post{post1, post2}).Error; err != nil {
		log.Fatalf("创建文章失败：%v", err)
	}
	// 验证用户文章数量（应变为 2）
	var updatedUser User
	db.First(&updatedUser, user.ID)
	log.Printf("用户文章数量更新后：%d", updatedUser.ArticleCount) // 输出：2

	// 创建评论（用于验证关联查询和钩子）
	comment1 := Comment{Content: "测试评论1", PostID: post1.ID}
	comment2 := Comment{Content: "测试评论2", PostID: post1.ID}
	comment3 := Comment{Content: "测试评论3", PostID: post2.ID}
	if err := db.Create(&[]Comment{comment1, comment2, comment3}).Error; err != nil {
		log.Fatalf("创建评论失败：%v", err)
	}

	// 关联查询演示（题目2）
	// 查询用户的文章及评论
	userWithArticles, err := getUserArticlesWithComments(db, user.ID)
	if err != nil {
		log.Printf("用户文章查询失败：%v", err)
	} else {
		log.Printf("\n=== 用户 %s 的文章及评论 ===", userWithArticles.Username)
		for _, post := range userWithArticles.Posts {
			log.Printf("文章：%s（评论数：%d）", post.Title, len(post.Comments))
			for _, comment := range post.Comments {
				log.Printf("  - 评论：%s", comment.Content)
			}
		}
	}

	// 查询评论最多的文章
	mostCommentedPost, err := getMostCommentedPost(db)
	if err != nil {
		log.Printf("评论最多文章查询失败：%v", err)
	} else {
		log.Printf("\n=== 评论最多的文章 ===")
		log.Printf("标题：%s，评论数：%d，状态：%s", mostCommentedPost.Title, mostCommentedPost.CommentCount, mostCommentedPost.CommentStatus)
	}

	// 钩子函数演示（题目3）
	// 删除 post1 的所有评论（触发 Comment 的 AfterDelete 钩子）
	if err := db.Where("post_id = ?", post1.ID).Delete(&Comment{}).Error; err != nil {
		log.Printf("删除评论失败：%v", err)
	}
	// 验证 post1 的状态（应变为 "无评论"）
	var updatedPost Post
	db.First(&updatedPost, post1.ID)
	log.Printf("\n=== 删除所有评论后文章状态 ===")
	log.Printf("文章：%s，状态：%s", updatedPost.Title, updatedPost.CommentStatus) // 输出：无评论
}

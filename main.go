package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type News struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Image       string `json:"image"`
	ArticleLink string `json:"articleLink"`
	Source      string `json:"source"`
	Author      string `json:"author"`
	PublishedAt string `json:"publishedAt"`
}

var db *sql.DB

func initDB() {
	var err error
	// 从环境变量获取数据库连接信息
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 设置默认值
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "root"
	}
	if dbPassword == "" {
		dbPassword = "123456"
	}
	if dbName == "" {
		dbName = "ai_news_db"
	}

	// 构建连接字符串（修改这一行）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建数据库（如果不存在）
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS ai_news_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal("创建数据库失败:", err)
	}

	// 使用数据库
	// 删除这两行USE语句（不再需要）
	// _, err = db.Exec("USE ai_news_db")
	// if err != nil {
	//     log.Fatal("使用数据库失败:", err)
	// }

	// 创建表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS news (
		id INT AUTO_INCREMENT PRIMARY KEY,
		category VARCHAR(100) NOT NULL,
		title VARCHAR(255) NOT NULL,
		summary TEXT NOT NULL,
		image VARCHAR(500) NOT NULL,
		article_link VARCHAR(500) NOT NULL,
		source VARCHAR(100),
		author VARCHAR(100),
		published_at VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_category (category),
		INDEX idx_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("创建表失败:", err)
	}

	// 检查表是否为空，如果是则插入示例数据
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
	if err != nil {
		log.Fatal("检查表数据失败:", err)
	}

	if count == 0 {
		insertSampleData()
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("数据库连接和表初始化成功")
}

func insertSampleData() {
	sampleData := []News{
		{
			Category:    "AI Simulation",
			Title:       "AI, Simulation, And The Generative Adversarial Network",
			Summary:     "AI enhances simulations by analyzing rich data, using tools like GANs, VAEs, and digital twins to significantly improve processes in manufacturing, healthcare, and synthetic data generation.",
			Image:       "https://images.unsplash.com/photo-1620712943543-bcc4688e7485?ixlib=rb-4.0.3&auto=format&fit=crop&w=1350&q=80",
			ArticleLink: "https://www.npr.org/2025/09/05/nx-s1-5529404/anthropic-settlement-authors-copyright-ai",
			Source:      "Tech Review",
			Author:      "AI Research Team",
			PublishedAt: "2024-01-15T10:30:00Z",
		},
		// ... 其他示例数据
	}

	stmt, err := db.Prepare(`INSERT INTO news (category, title, summary, image, article_link, source, author, published_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal("准备插入语句失败:", err)
	}
	defer stmt.Close()

	for _, news := range sampleData {
		_, err := stmt.Exec(news.Category, news.Title, news.Summary, news.Image, news.ArticleLink, news.Source, news.Author, news.PublishedAt)
		if err != nil {
			log.Printf("插入示例数据失败: %v", err)
		}
	}

	log.Println("示例数据插入完成")
}

func getNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 添加缓存控制头
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// 添加调试打印
	log.Println("收到GET /api/news请求")

	rows, err := db.Query("SELECT id, category, title, summary, image, article_link, source, author, published_at FROM news ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "数据库查询失败", http.StatusInternalServerError)
		log.Println("查询失败:", err)
		return
	}
	defer rows.Close()

	var newsList []News
	for rows.Next() {
		var news News
		err := rows.Scan(&news.ID, &news.Category, &news.Title, &news.Summary, &news.Image, &news.ArticleLink, &news.Source, &news.Author, &news.PublishedAt)
		if err != nil {
			http.Error(w, "数据读取失败", http.StatusInternalServerError)
			log.Println("数据读取失败:", err)
			return
		}
		newsList = append(newsList, news)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "数据遍历失败", http.StatusInternalServerError)
		log.Println("数据遍历失败:", err)
		return
	}

	// 添加查询结果打印
	log.Printf("查询到 %d 条新闻数据", len(newsList))
	for i, news := range newsList {
		log.Printf("新闻 %d: ID=%d, 标题=%s", i+1, news.ID, news.Title)
	}

	json.NewEncoder(w).Encode(newsList)
}

func getNewsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	var news News
	err = db.QueryRow("SELECT id, category, title, summary, image, article_link, source, author, published_at FROM news WHERE id = ?", id).Scan(
		&news.ID, &news.Category, &news.Title, &news.Summary, &news.Image, &news.ArticleLink, &news.Source, &news.Author, &news.PublishedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "数据库查询失败", http.StatusInternalServerError)
		log.Println("查询失败:", err)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func createNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var newsList []News
	err := json.NewDecoder(r.Body).Decode(&newsList)
	if err != nil {
		http.Error(w, "Invalid JSON format, expected array of news objects", http.StatusBadRequest)
		log.Println("JSON解析失败:", err)
		return
	}

	// 修改字段验证逻辑 - 跳过缺项数据而不是报错
	var validNewsList []News
	for i, news := range newsList {
		if news.Category == "" || news.Title == "" || news.Summary == "" || news.Image == "" || news.ArticleLink == "" {
			log.Printf("跳过第 %d 条新闻数据，缺少必填字段", i+1)
			continue // 跳过缺项数据
		}
		validNewsList = append(validNewsList, news)
	}

	if len(validNewsList) == 0 {
		http.Error(w, "没有有效的新闻数据", http.StatusBadRequest)
		return
	}

	// 先删除所有现有数据
	_, err = db.Exec("DELETE FROM news")
	if err != nil {
		http.Error(w, "删除现有数据失败", http.StatusInternalServerError)
		log.Println("删除数据失败:", err)
		return
	}
	log.Println("已删除所有现有新闻数据")

	stmt, err := db.Prepare(`INSERT INTO news (category, title, summary, image, article_link, source, author, published_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		http.Error(w, "数据库准备语句失败", http.StatusInternalServerError)
		log.Println("准备插入语句失败:", err)
		return
	}
	defer stmt.Close()

	var createdNews []News
	for _, news := range validNewsList {
		result, err := stmt.Exec(news.Category, news.Title, news.Summary, news.Image, news.ArticleLink, news.Source, news.Author, news.PublishedAt)
		if err != nil {
			log.Printf("插入新闻失败: %v，跳过此条数据", err)
			continue // 跳过插入失败的数据
		}

		id, _ := result.LastInsertId()
		news.ID = int(id)
		createdNews = append(createdNews, news)
		log.Printf("新闻创建成功: ID=%d, 标题=%s", id, news.Title)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNews)
	log.Printf("批量创建完成，共创建 %d 条有效新闻", len(createdNews))
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/news", getNews).Methods("GET")
	r.HandleFunc("/api/news", createNews).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/news/{id}", getNewsByID).Methods("GET")

	log.Println("Server is running on http://localhost:18081")
	log.Fatal(http.ListenAndServe(":18081", r))
}

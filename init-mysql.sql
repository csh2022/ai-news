CREATE DATABASE IF NOT EXISTS ai_news_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ai_news_db;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
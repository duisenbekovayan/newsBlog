package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var db *sql.DB

func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.127.126.24:3306)/golang_test")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS articles (id INT AUTO_INCREMENT, title VARCHAR(255),
        anons TEXT, full_text TEXT, PRIMARY KEY (id))`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE articles")
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}
}

func teardownTestDB() {
	db.Close()
}

func TestFetchArticles(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	_, err := db.Exec(`INSERT INTO articles (title, anons, full_text) VALUES 
        ('Test Title 1', 'Test Anons 1', 'Test Full Text 1'),
        ('Test Title 2', 'Test Anons 2', 'Test Full Text 2')`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		t.Fatalf("failed to fetch articles: %v", err)
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.Id, &article.Title, &article.Anons, &article.Fulltext)
		if err != nil {
			t.Fatalf("failed to scan article: %v", err)
		}
		articles = append(articles, article)
	}

	if len(articles) != 2 {
		t.Errorf("expected 2 articles, got %d", len(articles))
	}
}

func TestSaveArticle(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	title := "New Test Title"
	anons := "New Test Anons"
	fullText := "New Test Full Text"

	_, err := db.Exec("INSERT INTO articles (title, anons, full_text) VALUES (?, ?, ?)", title, anons, fullText)
	if err != nil {
		t.Fatalf("failed to insert article: %v", err)
	}

	var article Article
	err = db.QueryRow("SELECT id, title, anons, full_text FROM articles WHERE title = ?", title).Scan(&article.Id, &article.Title, &article.Anons, &article.Fulltext)
	if err != nil {
		t.Fatalf("failed to retrieve inserted article: %v", err)
	}

	if article.Title != title || article.Anons != anons || article.Fulltext != fullText {
		t.Errorf("retrieved article does not match inserted data: got %+v", article)
	}
}

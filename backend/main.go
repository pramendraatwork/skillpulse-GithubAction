package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Skill struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Level    string `json:"level"`
}

var db *sql.DB

func connectDB() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "skillpulse")
	password := getEnv("DB_PASSWORD", "skillpulse123")
	dbname := getEnv("DB_NAME", "skillpulsedb")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Warning: Could not connect to DB: %v. Using in-memory data.", err)
		db = nil
		return
	}

	if err = db.Ping(); err != nil {
		log.Printf("Warning: DB not reachable: %v. Using in-memory data.", err)
		db = nil
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS skills (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		category VARCHAR(100),
		level VARCHAR(50)
	)`)
	if err != nil {
		log.Printf("Warning: Could not create table: %v", err)
	}

	db.Exec(`INSERT INTO skills (name, category, level) VALUES
		('Docker', 'DevOps', 'Intermediate'),
		('Kubernetes', 'DevOps', 'Beginner'),
		('Terraform', 'IaC', 'Advanced'),
		('Go', 'Programming', 'Intermediate'),
		('GitHub Actions', 'CI/CD', 'Beginner')
		ON CONFLICT DO NOTHING`)

	log.Println("Database connected successfully!")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getSkills(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" { return }

	var skills []Skill

	if db != nil {
		rows, err := db.Query("SELECT id, name, category, level FROM skills")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var s Skill
				rows.Scan(&s.ID, &s.Name, &s.Category, &s.Level)
				skills = append(skills, s)
			}
		}
	}

	if len(skills) == 0 {
		skills = []Skill{
			{1, "Docker", "DevOps", "Intermediate"},
			{2, "Kubernetes", "DevOps", "Beginner"},
			{3, "Terraform", "IaC", "Advanced"},
			{4, "Go", "Programming", "Intermediate"},
			{5, "GitHub Actions", "CI/CD", "Beginner"},
		}
	}

	json.NewEncoder(w).Encode(skills)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "skillpulse-backend"})
}

func main() {
	connectDB()
	http.HandleFunc("/api/skills", getSkills)
	http.HandleFunc("/api/health", healthCheck)
	port := getEnv("PORT", "8080")
	log.Printf("SkillPulse Backend running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

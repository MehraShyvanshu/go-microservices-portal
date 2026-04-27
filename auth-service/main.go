package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var jwtKey = []byte("triconinfotech_secret_2024")
var db *sql.DB

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./auth.db")
	if err != nil {
		log.Fatal("❌ DB open failed:", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		department TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		email TEXT DEFAULT ''
	)`)
	if err != nil {
		log.Fatal("❌ Table creation failed:", err)
	}
	seedUsers()
}

func seedUsers() {
	seeds := []struct{ username, password, role, department, phone, email string }{
		{"admin", "admin123", "admin", "Management", "+91-9808181103", "admin@triconinfotech.com"},
		{"Shyvanshu Mehra", "Tricon@123", "admin", "Management", "+91-9808181103", "shyvanshu@triconinfotech.com"},
		{"john", "user123", "user", "Engineering", "+91-9452149721", "john@triconinfotech.com"},
		{"sarah", "user123", "hr", "HR", "+91-9876543210", "sarah@triconinfotech.com"},
		{"alex", "user123", "marketing", "Marketing", "+91-9123456789", "alex@triconinfotech.com"},
	}
	for _, s := range seeds {
		var count int
		db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", s.username).Scan(&count)
		if count == 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
			db.Exec("INSERT INTO users (username,password,role,department,phone,email) VALUES (?,?,?,?,?,?)",
				s.username, string(hash), s.role, s.department, s.phone, s.email)
			log.Printf("✅ Seeded user: %s (%s)", s.username, s.role)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)

	var user User
	var hashedPwd string
	err := db.QueryRow(
		"SELECT id,username,password,role,department,phone,email FROM users WHERE username=?",
		creds.Username,
	).Scan(&user.ID, &user.Username, &hashedPwd, &user.Role, &user.Department, &user.Phone, &user.Email)
	if err != nil {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(creds.Password)) != nil {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      token,
		"role":       user.Role,
		"username":   user.Username,
		"user_id":    user.ID,
		"department": user.Department,
		"email":      user.Email,
		"phone":      user.Phone,
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		Department string `json:"department"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&payload)
	if payload.Role == "" {
		payload.Role = "user"
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (username,password,role,department,phone,email) VALUES (?,?,?,?,?,?)",
		payload.Username, string(hash), payload.Role, payload.Department, payload.Phone, payload.Email)
	if err != nil {
		http.Error(w, `{"error":"Username already exists"}`, http.StatusConflict)
		return
	}
	w.Write([]byte(`{"success":true}`))
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id,username,role,department,phone,email FROM users")
	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Role, &u.Department, &u.Phone, &u.Email)
		users = append(users, u)
	}
	if users == nil {
		users = []User{}
	}
	json.NewEncoder(w).Encode(users)
}

func updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var p struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}
	json.NewDecoder(r.Body).Decode(&p)
	db.Exec("UPDATE users SET role=? WHERE id=?", p.Role, p.UserID)
	w.Write([]byte(`{"success":true}`))
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db.Exec("DELETE FROM users WHERE id=?", id)
	w.Write([]byte(`{"success":true}`))
}

func main() {
	initDB()
	http.HandleFunc("/login", cors(loginHandler))
	http.HandleFunc("/register", cors(registerHandler))
	http.HandleFunc("/users", cors(getUsersHandler))
	http.HandleFunc("/users/role", cors(updateRoleHandler))
	http.HandleFunc("/users/delete", cors(deleteUserHandler))

	fmt.Println("🔐 Auth Service running on :8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}

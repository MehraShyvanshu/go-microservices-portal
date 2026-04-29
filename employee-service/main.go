package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

var jwtKey = []byte("triconinfotech_secret_2024")
var db *sql.DB

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Employee struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Department  string `json:"department"`
	Phone       string `json:"phone"`
	Salary      string `json:"salary"`
	Status      string `json:"status"`
	JoinedDate  string `json:"joined_date"`
	AvatarColor string `json:"avatar_color"`
}

type Stats struct {
	Total    int            `json:"total"`
	Active   int            `json:"active"`
	Inactive int            `json:"inactive"`
	ByDept   map[string]int `json:"by_department"`
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

func authorize(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("no token")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(auth[7:], claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.Role, nil
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./employees.db")
	if err != nil {
		log.Fatal("❌ DB open failed:", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		department TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		salary TEXT DEFAULT '$0',
		status TEXT DEFAULT 'active',
		joined_date TEXT DEFAULT '',
		avatar_color TEXT DEFAULT '#6366f1'
	)`)
	if err != nil {
		log.Fatal("❌ Table creation failed:", err)
	}
	seedEmployees()
}

func seedEmployees() {
	seeds := []Employee{
		{Name: "John Smith", Email: "john@triconinfotech.com", Role: "user", Department: "Engineering", Phone: "+91-9452149721", Salary: "$95,000", Status: "active", JoinedDate: "2022-03-15", AvatarColor: "#6366f1"},
		{Name: "Sarah Johnson", Email: "sarah@triconinfotech.com", Role: "user", Department: "HR", Phone: "+91-9876543210", Salary: "$75,000", Status: "active", JoinedDate: "2021-07-20", AvatarColor: "#ec4899"},
		{Name: "Alex Chen", Email: "alex@triconinfotech.com", Role: "user", Department: "Marketing", Phone: "+91-9123456789", Salary: "$85,000", Status: "active", JoinedDate: "2023-01-10", AvatarColor: "#f59e0b"},
		{Name: "Emily Davis", Email: "emily@triconinfotech.com", Role: "user", Department: "Finance", Phone: "+91-9452149721", Salary: "$90,000", Status: "active", JoinedDate: "2020-11-05", AvatarColor: "#10b981"},
		{Name: "Michael Brown", Email: "michael@triconinfotech.com", Role: "admin", Department: "Engineering", Phone: "+91-9876543210", Salary: "$120,000", Status: "active", JoinedDate: "2019-06-01", AvatarColor: "#3b82f6"},
		{Name: "Lisa Wilson", Email: "lisa@triconinfotech.com", Role: "user", Department: "Operations", Phone: "+91-9123456789", Salary: "$70,000", Status: "inactive", JoinedDate: "2022-09-15", AvatarColor: "#8b5cf6"},
		{Name: "David Martinez", Email: "david@triconinfotech.com", Role: "user", Department: "Engineering", Phone: "+91-9452149721", Salary: "$100,000", Status: "active", JoinedDate: "2021-04-20", AvatarColor: "#14b8a6"},
		{Name: "Jessica Lee", Email: "jessica@triconinfotech.com", Role: "user", Department: "Marketing", Phone: "+91-9876543210", Salary: "$80,000", Status: "active", JoinedDate: "2023-03-01", AvatarColor: "#f97316"},
		{Name: "Shyvanshu Mehra", Email: "shyvanshu@triconinfotech.com", Role: "admin", Department: "Engineering", Phone: "+91-9452149721", Salary: "$120,000", Status: "active", JoinedDate: "2019-06-01", AvatarColor: "#042b69ff"},
	}
	for _, e := range seeds {
		var count int
		db.QueryRow("SELECT COUNT(*) FROM employees WHERE email=?", e.Email).Scan(&count)
		if count == 0 {
			db.Exec(`INSERT INTO employees (name,email,role,department,phone,salary,status,joined_date,avatar_color)
				VALUES (?,?,?,?,?,?,?,?,?)`,
				e.Name, e.Email, e.Role, e.Department, e.Phone, e.Salary, e.Status, e.JoinedDate, e.AvatarColor)
			log.Printf("✅ Seeded employee: %s", e.Name)
		}
	}
}

func scanEmployee(rows *sql.Rows) (Employee, error) {
	var e Employee
	err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Department, &e.Phone, &e.Salary, &e.Status, &e.JoinedDate, &e.AvatarColor)
	return e, err
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id,name,email,role,department,phone,salary,status,joined_date,avatar_color FROM employees ORDER BY id ASC")
	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	employees := []Employee{}
	for rows.Next() {
		e, _ := scanEmployee(rows)
		employees = append(employees, e)
	}
	json.NewEncoder(w).Encode(employees)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	role, err := authorize(r)
	if err != nil || role != "admin" {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}
	var e Employee
	json.NewDecoder(r.Body).Decode(&e)
	if e.JoinedDate == "" {
		e.JoinedDate = time.Now().Format("2006-01-02")
	}
	if e.AvatarColor == "" {
		colors := []string{"#6366f1", "#ec4899", "#f59e0b", "#10b981", "#3b82f6", "#8b5cf6", "#14b8a6", "#f97316"}
		e.AvatarColor = colors[len(e.Name)%len(colors)]
	}
	if e.Status == "" {
		e.Status = "active"
	}
	if e.Role == "" {
		e.Role = "user"
	}

	result, err := db.Exec(`INSERT INTO employees (name,email,role,department,phone,salary,status,joined_date,avatar_color)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		e.Name, e.Email, e.Role, e.Department, e.Phone, e.Salary, e.Status, e.JoinedDate, e.AvatarColor)
	if err != nil {
		http.Error(w, `{"error":"Email already exists or DB error"}`, http.StatusConflict)
		return
	}
	id, _ := result.LastInsertId()
	e.ID = int(id)
	json.NewEncoder(w).Encode(e)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	role, err := authorize(r)
	if err != nil || role != "admin" {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	var e Employee
	json.NewDecoder(r.Body).Decode(&e)
	_, err = db.Exec(`UPDATE employees SET name=?,email=?,role=?,department=?,phone=?,salary=?,status=? WHERE id=?`,
		e.Name, e.Email, e.Role, e.Department, e.Phone, e.Salary, e.Status, id)
	if err != nil {
		log.Printf("❌ Update failed for employee %d: %v", id, err)
		http.Error(w, `{"error":"Update failed"}`, http.StatusInternalServerError)
		return
	}
	log.Printf("✅ Updated employee %d: %s", id, e.Name)
	e.ID = id
	json.NewEncoder(w).Encode(e)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	role, err := authorize(r)
	if err != nil || role != "admin" {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}
	idStr := r.URL.Query().Get("id")
	db.Exec("DELETE FROM employees WHERE id=?", idStr)
	w.Write([]byte(`{"success":true}`))
}

func getStats(w http.ResponseWriter, r *http.Request) {
	var stats Stats
	stats.ByDept = make(map[string]int)
	db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&stats.Total)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE status='active'").Scan(&stats.Active)
	db.QueryRow("SELECT COUNT(*) FROM employees WHERE status='inactive'").Scan(&stats.Inactive)
	rows, _ := db.Query("SELECT department, COUNT(*) FROM employees GROUP BY department")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var dept string
			var cnt int
			rows.Scan(&dept, &cnt)
			stats.ByDept[dept] = cnt
		}
	}
	json.NewEncoder(w).Encode(stats)
}

func main() {
	initDB()
	http.HandleFunc("/employees", cors(getEmployees))
	http.HandleFunc("/employees/create", cors(createEmployee))
	http.HandleFunc("/employees/update", cors(updateEmployee))
	http.HandleFunc("/employees/delete", cors(deleteEmployee))
	http.HandleFunc("/employees/stats", cors(getStats))

	fmt.Println("👥 Employee Service running on :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}

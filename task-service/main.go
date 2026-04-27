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

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  string `json:"assigned_to"`
	CreatedBy   string `json:"created_by"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
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

func getClaimsFromRequest(r *http.Request) (*Claims, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, fmt.Errorf("no token")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(auth[7:], claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal("❌ DB open failed:", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		status TEXT NOT NULL DEFAULT 'todo',
		priority TEXT NOT NULL DEFAULT 'medium',
		assigned_to TEXT DEFAULT '',
		created_by TEXT DEFAULT '',
		due_date TEXT DEFAULT '',
		created_at TEXT DEFAULT ''
	)`)
	if err != nil {
		log.Fatal("❌ Table creation failed:", err)
	}
	seedTasks()
}

func seedTasks() {
	seeds := []Task{
		{Title: "Design new landing page", Description: "Create wireframes and mockups for the revamped marketing website", Status: "done", Priority: "high", AssignedTo: "alex", CreatedBy: "admin", DueDate: "2024-04-10"},
		{Title: "Backend API optimization", Description: "Optimize slow database queries in the employee service", Status: "in-progress", Priority: "high", AssignedTo: "john", CreatedBy: "admin", DueDate: "2024-04-22"},
		{Title: "Q2 HR Policy Updates", Description: "Review and update remote work and leave policies for Q2", Status: "in-progress", Priority: "medium", AssignedTo: "sarah", CreatedBy: "admin", DueDate: "2024-04-25"},
		{Title: "Employee onboarding flow", Description: "Streamline the onboarding checklist for new hires", Status: "todo", Priority: "medium", AssignedTo: "sarah", CreatedBy: "admin", DueDate: "2024-04-30"},
		{Title: "Fix production bug #204", Description: "Resolve the login timeout issue reported by multiple users", Status: "done", Priority: "high", AssignedTo: "john", CreatedBy: "admin", DueDate: "2024-04-12"},
		{Title: "Marketing campaign analytics", Description: "Set up tracking for the April email campaign", Status: "todo", Priority: "low", AssignedTo: "alex", CreatedBy: "admin", DueDate: "2024-05-01"},
		{Title: "Security audit", Description: "Conduct quarterly security review of all microservices", Status: "todo", Priority: "high", AssignedTo: "john", CreatedBy: "admin", DueDate: "2024-04-28"},
		{Title: "Update employee documentation", Description: "Refresh internal wiki with current team structure and processes", Status: "in-progress", Priority: "low", AssignedTo: "sarah", CreatedBy: "admin", DueDate: "2024-05-05"},
	}
	var count int
	db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if count == 0 {
		for _, t := range seeds {
			t.CreatedAt = time.Now().Format("2006-01-02")
			db.Exec(`INSERT INTO tasks (title,description,status,priority,assigned_to,created_by,due_date,created_at)
				VALUES (?,?,?,?,?,?,?,?)`,
				t.Title, t.Description, t.Status, t.Priority, t.AssignedTo, t.CreatedBy, t.DueDate, t.CreatedAt)
		}
		log.Printf("✅ Seeded %d tasks", len(seeds))
	}
}

func routeTasks(w http.ResponseWriter, r *http.Request) {
	// Route: /tasks or /tasks/:id
	path := strings.TrimPrefix(r.URL.Path, "/tasks")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			listTasks(w, r)
		case http.MethodPost:
			addTask(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
		return
	}

	// Has ID
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		editTask(w, r, id)
	case http.MethodDelete:
		removeTask(w, r, id)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id,title,description,status,priority,assigned_to,created_by,due_date,created_at FROM tasks ORDER BY id DESC")
	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.CreatedBy, &t.DueDate, &t.CreatedAt)
		tasks = append(tasks, t)
	}
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	var t Task
	json.NewDecoder(r.Body).Decode(&t)
	if t.Status == "" {
		t.Status = "todo"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	t.CreatedBy = claims.Username
	t.CreatedAt = time.Now().Format("2006-01-02")

	result, err := db.Exec(`INSERT INTO tasks (title,description,status,priority,assigned_to,created_by,due_date,created_at)
		VALUES (?,?,?,?,?,?,?,?)`,
		t.Title, t.Description, t.Status, t.Priority, t.AssignedTo, t.CreatedBy, t.DueDate, t.CreatedAt)
	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	t.ID = int(id)
	json.NewEncoder(w).Encode(t)
}

func editTask(w http.ResponseWriter, r *http.Request, id int) {
	_, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	var t Task
	json.NewDecoder(r.Body).Decode(&t)
	db.Exec(`UPDATE tasks SET title=?,description=?,status=?,priority=?,assigned_to=?,due_date=? WHERE id=?`,
		t.Title, t.Description, t.Status, t.Priority, t.AssignedTo, t.DueDate, id)
	t.ID = id
	json.NewEncoder(w).Encode(t)
}

func removeTask(w http.ResponseWriter, r *http.Request, id int) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	_ = claims
	db.Exec("DELETE FROM tasks WHERE id=?", id)
	w.Write([]byte(`{"success":true}`))
}

func main() {
	initDB()
	http.HandleFunc("/tasks", cors(func(w http.ResponseWriter, r *http.Request) {
		routeTasks(w, r)
	}))
	http.HandleFunc("/tasks/", cors(func(w http.ResponseWriter, r *http.Request) {
		routeTasks(w, r)
	}))

	fmt.Println("📋 Task Service running on :8003")
	log.Fatal(http.ListenAndServe(":8003", nil))
}

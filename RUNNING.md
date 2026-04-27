# 🚀 How to Run Tricon Infotech Employee Portal

## Prerequisites

| Tool | Version | Download |
|------|---------|----------|
| Go | 1.21+ | https://go.dev/dl/ |
| Node.js | 18+ | https://nodejs.org |
| npm | 8+ | (bundled with Node.js) |
| DBeaver CE | Any | https://dbeaver.io/download/ (optional, to view DB) |

---

## Step 1 — Clone / Open the Project

Make sure you're in the project root:
```
cd c:\Users\Shyvanshu.Mehra\go-microservices-portal
```

---

## Step 2 — Start the Auth Service (Terminal 1)

```powershell
cd auth-service
go run main.go
```

**Expected output:**
```
✅ Seeded user: admin (admin)
✅ Seeded user: john (user)
✅ Seeded user: sarah (user)
✅ Seeded user: alex (user)
🔐 Auth Service running on :8001
```

The file `auth.db` will be created automatically in the `auth-service/` folder.

---

## Step 3 — Start the Employee Service (Terminal 2)

```powershell
cd employee-service
go run main.go
```

**Expected output:**
```
✅ Seeded employee: John Smith
... (8 employees seeded)
👥 Employee Service running on :8002
```

The file `employees.db` will be created automatically in the `employee-service/` folder.

---

## Step 4 — Start the Task Service (Terminal 3)

```powershell
cd task-service
go run main.go
```

**Expected output:**
```
✅ Seeded 8 tasks
📋 Task Service running on :8003
```

The file `tasks.db` will be created automatically in the `task-service/` folder.

---

## Step 5 — Start the Frontend (Terminal 4)

```powershell
cd frontend
npm start
```

The React app will automatically open at **http://localhost:3000**

---

## Step 6 — Login and Explore

| Demo Account | Username | Password | Access Level |
|---|---|---|---|
| Admin | `admin` | `admin123` | Full access |
| Engineer | `john`  | `user123`  | Standard access |
| HR | `sarah` | `user123`  | Standard access |
| Marketing | `alex`  | `user123`  | Standard access |

> 💡 **Tip:** On the login page, click any demo credential button to auto-fill the form!

---

## Step 7 — View Data in DBeaver (Optional)

1. Open **DBeaver Community**
2. Click the plug icon → **New Database Connection**
3. Search and select **SQLite**
4. Click **Browse** and navigate to:
   - `c:\Users\Shyvanshu.Mehra\go-microservices-portal\auth-service\auth.db`
5. Click **Test Connection** → **Finish**
6. Expand **Tables** to see `users` table with all data
7. Repeat for `employee-service\employees.db` and `task-service\tasks.db`

---

## Ports Summary

| Service | Port | URL |
|---------|------|-----|
| Frontend (React) | 3000 | http://localhost:3000 |
| Auth Service | 8001 | http://localhost:8001 |
| Employee Service | 8002 | http://localhost:8002 |
| Task Service | 8003 | http://localhost:8003 |

---

## Troubleshooting

### Go services fail to start
- Make sure Go 1.21+ is installed: `go version`
- Run `go mod tidy` in the failing service directory

### npm start fails
- Run `npm install` in the `frontend/` directory
- Make sure Node.js 18+ is installed: `node --version`

### Login shows error / charts empty
- Make sure the auth service is running on port 8001
- Employee/task data requires services on 8002 and 8003
- Check terminal windows for any error messages

### Database already has stale data
- Delete the `.db` file in the service directory
- Restart the service — it will re-seed fresh data

---

## Optional: API Gateway

The `api-gateway` (optional Gin proxy) can be started if needed:
```powershell
cd api-gateway
go run .
```
It runs on port **8080** and proxies to the microservices.  
The frontend currently connects **directly** to each service for simplicity.

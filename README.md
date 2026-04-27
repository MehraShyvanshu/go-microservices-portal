# Tricon Infotech Employee Portal

<div align="center">

![Tricon Infotech](https://img.shields.io/badge/Tricon Infotech-Employee%20Portal-7c3aed?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react)
![SQLite](https://img.shields.io/badge/SQLite-DBeaver-003B57?style=for-the-badge&logo=sqlite)

**Empowering People · Driving Growth**

</div>

---

## 🚀 Overview

Tricon Infotech Employee Portal is a **full-stack microservices application** built with Go (backend) and React (frontend). It provides a comprehensive employee management system with role-based access control, task management, announcements, leave requests, and more — all with a stunning, animation-rich UI.

---

## ✨ Features

### For Everyone
- 🎬 **Animated Splash Screen** — Company logo reveal with smooth framer-motion animations
- 🔐 **JWT Authentication** — Secure login with bcrypt-hashed passwords
- 📊 **Role-Based Dashboard** — Personalized dashboards for admin and user roles
- 👥 **Employee Directory** — Browse the full team with avatars, departments, and contact info
- 📋 **Kanban Task Board** — Drag-and-drop task management across Todo/In Progress/Done
- 📢 **Announcements** — Company-wide news and updates feed
- 🗓️ **Leave Requests** — Submit and track your leave applications
- 👤 **Profile Page** — Personal info, leave summary, and access rights matrix
- 🌙 **Dark/Light Mode** — Theme toggle, persisted across sessions

### Admin Only
- ➕ **Create Employees** — Full CRUD with department, salary, phone, role assignment
- ✏️ **Edit/Delete Employees** — Manage your entire workforce
- ✅ **Approve/Reject Leaves** — Review all pending leave requests
- 📝 **Post Announcements** — Broadcast company-wide messages with priority levels
- ⚙️ **Settings** — Add new users, change roles, view system info

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                   React Frontend                     │
│              (Port 3000 — npm start)                 │
└──────────────────┬──────────────────────────────────┘
                   │ HTTP/REST
      ┌────────────┼────────────┐
      ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│  Auth    │ │ Employee │ │   Task   │
│ Service  │ │ Service  │ │ Service  │
│  :8001   │ │  :8002   │ │  :8003   │
└────┬─────┘ └────┬─────┘ └────┬─────┘
     │             │             │
┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐
│ auth.db  │ │employees │ │ tasks.db │
│ (SQLite) │ │   .db    │ │ (SQLite) │
└──────────┘ └──────────┘ └──────────┘
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend Framework** | React 19 |
| **Animations** | Framer Motion |
| **Charts** | Recharts |
| **Icons** | React Icons + Emoji |
| **State Management** | Zustand (with persist) |
| **Routing** | React Router DOM v7 |
| **Backend Language** | Go 1.21+ |
| **Database** | SQLite via `modernc.org/sqlite` (no CGO) |
| **Authentication** | JWT (golang-jwt/jwt/v5) + bcrypt |
| **API Style** | RESTful HTTP |
| **DB Viewer** | DBeaver Community (open `.db` files) |

---

## 👤 Demo Credentials

| Username | Password | Role | Department |
|---------|---------|------|-----------|
| `admin` | `admin123` | 👑 Admin | Management |
| `john` | `user123` | 👤 User | Engineering |
| `sarah` | `user123` | 👤 User | HR |
| `alex` | `user123` | 👤 User | Marketing |

---

## 📁 Project Structure

```
go-microservices-portal/
├── auth-service/          # JWT auth + user management
│   ├── main.go
│   ├── auth.db            # auto-created on first run
│   └── go.mod
├── employee-service/      # Employee CRUD
│   ├── main.go
│   ├── employees.db       # auto-created on first run
│   └── go.mod
├── task-service/          # Kanban task management
│   ├── main.go
│   ├── tasks.db           # auto-created on first run
│   └── go.mod
├── api-gateway/           # Optional Gin proxy
│   └── main.go
├── frontend/              # React application
│   ├── public/index.html
│   └── src/
│       ├── pages/         # Splash, Login, Dashboard, Employees, Board, ...
│       ├── components/    # Layout, Sidebar, Topbar, ProtectedRoute
│       ├── services/      # authService.js, api.js
│       └── store/         # useStore.js (Zustand)
├── FEATURES.md
├── RUNNING.md
└── README.md
```

---

## 🗃️ Viewing the Database in DBeaver

1. Download and open **DBeaver Community** (free)
2. Click **New Database Connection** → choose **SQLite**
3. Browse to the `.db` file inside the service directory:
   - `auth-service/auth.db` → users table
   - `employee-service/employees.db` → employees table
   - `task-service/tasks.db` → tasks table
4. Click **Finish** and expand the tables to see all your data live!

---

## 📄 Documentation

- 📋 [FEATURES.md](./FEATURES.md) — Detailed feature list, code flow, and architecture
- 🚀 [RUNNING.md](./RUNNING.md) — Step-by-step run instructions

# Tricon Infotech Employee Portal — Features & Code Flow

## 🌟 Feature List

### Authentication & Security
- **JWT-based authentication** — Tokens include username, role, and user_id
- **bcrypt password hashing** — Industry-standard hashing (cost 10)
- **Role-based access control** — Admin and User roles with different permissions
- **Protected routes** — Frontend route guards using React Router DOM
- **Token persistence** — JWT stored in localStorage, cleared on logout

### Employee Management (Admin)
- **Full CRUD** — Create, Read, Update, Delete employees
- **Rich employee profiles** — Name, email, department, phone, salary, status, join date, role
- **Avatar colors** — Auto-assigned per employee name for visual identification
- **Search & filter** — Real-time search by name/email + filter by department/status
- **Animated cards** — Framer Motion entrance animations with hover effects

### Task Board (Kanban)
- **Drag & Drop** — Native HTML5 drag-and-drop to move tasks between columns
- **3 columns** — Todo → In Progress → Done
- **Priority badges** — Color-coded: High (red), Medium (yellow), Low (green)
- **Assignee avatars** — Colored initials showing who owns each task
- **Due dates** — Track deadlines on each task card
- **Optimistic updates** — UI updates immediately, syncs with API in background

### Announcements
- **Admin post form** — Title, body, priority level (high/medium/low)
- **Priority filter** — Quick filter buttons to see announcements by priority
- **Color-coded left borders** — Visual priority indicators on each card
- **Animation** — Slide-in animations for each announcement

### Leave Management
- **User application form** — Type, from date, to date, reason
- **Admin review table** — Approve or Reject pending requests
- **Status tracking** — Pending / Approved / Rejected badges
- **Summary stats** — Mini stat cards showing counts per status
- **Scoped view** — Users see only their own; admins see all

### Dashboard
- **Animated stat cards** — Numbers count up from 0 using custom hook
- **Bar chart** — Employees grouped by department (Recharts)
- **Recent employees** — Latest team members with avatar + dept
- **Activity feed** — Recent portal activity log
- **Announcements preview** — Latest 3 announcements inline
- **Quick actions** — Admin shortcut buttons to key pages

### Profile
- **Personal info card** — All user details from JWT
- **Leave summary** — Total/Approved/Pending counts
- **Access rights matrix** — Shows allowed/restricted features by role

### Settings (Admin)
- **User management** — List all users, change roles inline, delete users
- **Add new user form** — Create users with username/password/role/department
- **Theme toggle** — Persistent dark/light mode
- **System info panel** — Tech stack details

---

## 🔄 Code Flow

### Authentication Flow
```
User types credentials → Login.js
  → authService.loginUser(username, password)
    → POST http://localhost:8001/login
      → auth-service/main.go:loginHandler
        → Query SQLite: SELECT * FROM users WHERE username=?
        → bcrypt.CompareHashAndPassword
        → jwt.NewWithClaims (includes username, role, user_id)
        → Returns { token, role, username, department, email }
  → Store in localStorage
  → navigate('/dashboard')
```

### Protected Route Flow
```
Navigate to /employees
  → ProtectedRoute checks localStorage.getItem('token')
  → If no token → redirect to /login
  → If token exists → check role vs allowed roles
  → If role mismatch → redirect to /unauthorized
  → Otherwise → render <Employees />
```

### Employee CRUD Flow (Admin)
```
Admin clicks "Add Employee"
  → EmployeeModal opens
  → API call: createEmployee(form)
    → POST http://localhost:8002/employees/create
      → employee-service/main.go:createEmployee
        → Checks JWT Authorization (role must be "admin")
        → INSERT INTO employees
        → Returns new employee with auto-generated ID
  → UI updates with AnimatePresence (card slides in)
```

### Kanban Drag-and-Drop Flow
```
User drags task card
  → onDragStart stores task.id in dataTransfer
User drops on column
  → onDrop reads task.id, finds task in state
  → Immediately updates task.status in local state (optimistic)
  → API call: updateTask(id, { ...task, status: newStatus })
    → PUT http://localhost:8003/tasks/:id
      → task-service routes to editTask()
      → UPDATE tasks SET status=? WHERE id=?
```

### Zustand Store (Announcements & Leaves)
```
useStore.js (with persist middleware)
  → Persists: isDark, announcements, leaves to localStorage key 'triconinfotech-portal'
  → On addAnnouncement: prepends new item to array
  → On updateLeaveStatus: maps over array, updates matching ID
  → Theme toggle: adds/removes 'light' class on document.body
```

---

## 🏛️ Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                    REACT FRONTEND (Port 3000)                 │
│                                                              │
│  App.js                                                      │
│  ├── Splash.js         (animated intro, auto-dismiss)        │
│  ├── Login.js          (real credentials form)               │
│  ├── ProtectedRoute    (RBAC guard)                          │
│  │                                                           │
│  └── Layout (Sidebar + Topbar)                               │
│      ├── Dashboard.js  (stats, charts, feed)                 │
│      ├── Employees.js  (CRUD, search, grid)                  │
│      ├── Board.js      (Kanban drag-and-drop)                │
│      ├── Announcements.js (post/filter/delete)               │
│      ├── Leaves.js     (apply/approve/reject)                │
│      ├── Profile.js    (info, rights matrix)                 │
│      └── Settings.js   (users, theme, sys info)              │
│                                                              │
│  State: useStore (Zustand) │ Services: authService, api.js   │
└──────────┬───────────────────────────────┬──────────────────┘
           │                               │
    /login, /register              /employees, /tasks
           │                               │
┌──────────▼──────────┐     ┌─────────────▼──────┐     ┌──────────────────┐
│   Auth Service      │     │  Employee Service   │     │  Task Service    │
│   Go  :8001         │     │  Go  :8002          │     │  Go  :8003       │
│                     │     │                     │     │                  │
│  POST /login        │     │  GET  /employees    │     │  GET  /tasks     │
│  POST /register     │     │  POST /create       │     │  POST /tasks     │
│  GET  /users        │     │  PUT  /update       │     │  PUT  /tasks/:id │
│  PUT  /users/role   │     │  DELETE /delete     │     │  DELETE /tasks/id│
│  DELETE /users/del  │     │  GET  /stats        │     │                  │
└──────────┬──────────┘     └──────────┬──────────┘     └───────┬──────────┘
           │                           │                         │
    ┌──────▼──────┐             ┌──────▼──────┐          ┌──────▼──────┐
    │  auth.db    │             │employees.db │          │  tasks.db   │
    │  (SQLite)   │             │  (SQLite)   │          │  (SQLite)   │
    │  users      │             │  employees  │          │  tasks      │
    └─────────────┘             └─────────────┘          └─────────────┘
           │                           │                         │
    ┌──────▼─── DBeaver ──────────────▼─────────────────────────▼──────┐
    │  Open .db files directly in DBeaver Community to view all data    │
    └───────────────────────────────────────────────────────────────────┘
```

---

## 📦 Dependencies

### Frontend
| Package | Version | Purpose |
|---------|---------|---------|
| react | 19 | UI framework |
| react-router-dom | 7 | Client-side routing |
| framer-motion | latest | Animations |
| recharts | latest | Data visualization charts |
| react-icons | latest | Icon library |
| zustand | 5 | Global state management |
| axios | 1 | HTTP client (available) |

### Backend (per service)
| Package | Purpose |
|---------|---------|
| golang-jwt/jwt/v5 | JWT creation and validation |
| modernc.org/sqlite | Pure-Go SQLite (no CGO needed) |
| golang.org/x/crypto | bcrypt password hashing |
| github.com/gin-gonic/gin | API Gateway (gin-based proxy) |
| github.com/gorilla/websocket | (legacy, replaced by REST) |

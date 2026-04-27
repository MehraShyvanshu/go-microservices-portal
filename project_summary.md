# Go Microservices Portal - Project Summary

## 1. Project Overview
The **Go Microservices Portal** is a full-stack web application implementing a microservices architecture. It demonstrates a distributed system on the backend written in Go and a single-page React application on the frontend. The entire application is containerized and orchestrated using **Docker Compose**.

---

## 2. Backend Stack & Architecture

### **Technologies Applied:**
- **Language**: Go (Golang)
- **Web Framework**: Gin (Gin-Gonic)
- **Authentication**: JWT (`golang-jwt`) for stateless token-based auth
- **Security**: `bcrypt` for secure password hashing
- **Containerization**: Docker

### **Microservices Architecture:**
The backend is divided into focused, specialized services communicating behind an API Gateway:

1. **API Gateway (`api-gateway` - Port 8080):**
   - Serves as the central entry point for the frontend, proxying requests to the appropriate internal services.
   - **Middlewares implemented**:
     - *Authentication Middleware*: Validates user JWT tokens.
     - *Role-Based Access Control (RBAC)*: Restricts endpoint access depending on whether the token role is `admin` or `user`.
2. **Auth Service (`auth-service` - Port 8001):**
   - Handles public endpoints: `POST /register` and `POST /login`.
   - Responsible for hashing passwords, issuing JWT tokens, and providing an endpoint (`PUT /role`) to manage user roles.
3. **User Service (`user-service` - Port 8002):**
   - Handles user data retrieval via the `GET /users` endpoint (proxy available only to admins).
4. **Task Service (`task-service` - Port 8003):**
   - Manages tasks via `GET /tasks` and `POST /tasks` endpoints (accessible by both admins and users).

---

## 3. Frontend Stack & Architecture

### **Technologies Applied:**
- **Library**: React (v19)
- **State Management**: Zustand (for global application state, handling tasks)
- **Routing**: React Router (`react-router-dom`)
- **HTTP Client**: Axios (for communicating with the API Gateway)
- **Styling**: Vanilla CSS (e.g., card UI elements, flexbox layouts)
- **Build Tool**: Create React App (react-scripts)

---

## 4. Frontend Features & Capabilities

Here are the features actively implemented in the frontend application that users can perform:

### **1. Authentication (Login Page)**
- Accepts an email and password to log in.
- Authenticates the user and redirects them to the main board upon successful validation.

### **2. Kanban Task Board (Board Page - `/board`)**
- Displays tasks fetched from the backend organized into interactive columns: `todo`, `in-progress`, and `done`.
- **Drag-and-Drop functionality**: Users can pick up a task card and move it across columns to change its status.
- Updates are instantly stored in the global Zustand state and dispatched to the backend via a background `POST /api/tasks` request, keeping data in sync.

### **3. Dashboard Analytics (Dashboard Page)**
- Fetches real-time task data from the backend to present analytics.
- Displays summary cards for total task counts.
- Renders an **Activity Chart** visualizing the volume or progression of tasks.

### **4. User & Role Management (Users Page - `/users`)**
- Fetches the active list of registered users.
- Gives administrative controls over user permissions.
- You can upgrade a user's permissions by clicking "**Admin**" or demote them by clicking "**User**", which immediately triggers a role update API call to the backend.

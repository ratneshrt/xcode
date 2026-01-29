# 🚀 Xcode – Online Judge Backend (LeetCode-like Platform)

Xcode is a **backend system for an online coding platform**, inspired by LeetCode and Codeforces, built using **Go (Golang)**.  
It supports secure authentication, problem management, asynchronous code execution, and sandboxed runtime environments using Docker.

This project focuses on **real-world backend engineering**, **security**, and **scalability**, rather than only solving algorithmic problems.

---

## ✨ Features

- 🔐 JWT-based authentication
- 👤 Role-based access control (admin / user)
- 📚 Problem management using YAML (Content as Code)
- 🧪 Support for public and hidden test cases
- ⚙️ Asynchronous submission processing
- 📦 Redis-backed job queue
- 🐳 Secure Docker-based code execution
- 🧱 Non-root sandboxed execution
- 🔁 Idempotent background workers
- 🧠 Production-grade system design

---

## 🧩 System Components

### 1️⃣ API Server

**Technology Stack**
- Go (Golang)
- Gin framework
- PostgreSQL
- JWT
- bcrypt

**Responsibilities**
- User registration and login
- JWT generation and validation
- Role-based authorization
- Problem listing and retrieval
- Accepting code submissions
- Persisting submission metadata

User-submitted code is **never executed directly** by the API server.

---

### 2️⃣ Problem Worker (Problem Ingestion Service)

Problems are managed as **YAML files**, not manual database entries.

**Responsibilities**
- Scans all `problem.yaml` files
- Parses YAML into Go structs
- Validates schema and constraints
- Inserts or updates problems in the database
- Replaces examples and test cases transactionally

**Key design points**
- Idempotent execution (safe to run multiple times)
- Uses database transactions to prevent partial writes
- Supports `draft` and `published` problem states
- Separates examples (human-readable) from test cases (machine-executed)

This approach follows **Content as Code**, commonly used in production systems.

---

### 3️⃣ Submission Processing

When a user submits code:
- The submission is stored with a `pending` status
- A job is pushed to Redis
- The API responds immediately

This design ensures low latency and allows the system to scale independently of execution workloads.

---

### 4️⃣ Execution Worker (Judge Engine)


**Responsibilities**
- Consumes submission jobs from Redis
- Fetches submission data and problem test cases
- Executes code in a Docker sandbox
- Determines the final verdict
- Updates submission status and result in the database

**Supported verdicts**
- `accepted`
- `wrong_answer`
- `runtime_error`
- `time_limit_exceeded`

---

## 🐳 Secure Code Execution with Docker

User code is considered **untrusted** and is always executed in isolation.

### Docker Runner Image

A custom Docker image is used as a controlled execution environment.

**Key security measures**
- Non-root user inside the container
- Isolated working directory (`/app`)
- No direct execution on the host
- Fresh container per submission
- Container destroyed after execution

This design closely resembles the execution model used by real online judges.

---


---

## ⚙️ Setup Requirements

- Go 1.22+
- Docker / Docker Desktop
- PostgreSQL
- Redis

---

## 🧠 Key Engineering Concepts Demonstrated

- Distributed backend design
- Asynchronous job processing
- Secure sandboxed execution
- Docker container lifecycle management
- Go module and package organization
- Database transactions and consistency
- Idempotent background workers
- Production debugging and fault isolation

---

## 🎯 Interview Talking Points

You can confidently describe this project as:

> “A LeetCode-like backend platform built in Go, featuring asynchronous code execution using Redis and worker processes, secure Docker-based sandboxes for untrusted code, and YAML-based problem ingestion following a Content-as-Code approach.”

---

## 🚀 Future Enhancements

- Support for additional languages (Python, Java)
- Resource usage tracking (time and memory)
- Firecracker-based microVM execution
- Execution metrics and observability
- Web frontend
- CI validation for problem YAML files
---
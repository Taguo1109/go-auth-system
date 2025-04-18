# 🛡️ go-auth-system

基於 Golang + Gin + JWT 的簡易用戶認證系統，包含註冊、登入、Token 刷新、登出功能，支援 Redis 快取與 Swagger API 文件。可作為開發實戰的模板或學習範例。

---

## 📦 技術棧

- [x] Golang 1.20+
- [x] Gin Web Framework
- [x] GORM（MySQL ORM）
- [x] JWT（使用 `github.com/golang-jwt/jwt/v5`）
- [x] Redis（使用 `go-redis` 套件）
- [x] Swagger 文檔（使用 `swaggo/swag`）
- [x] Kubernetes（使用 KIND 本地模擬）
- [x] Docker 支援（建置映像部署）

---

## 🚀 快速開始（本機執行）

### 1️⃣ 安裝依賴

```bash
go mod tidy
```

### 2️⃣ 初始化 Swagger 文檔

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
對 bash 使用者（Linux/macOS）：
```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```
```bash
swag init
```


### 3️⃣ 建立 .env.local 檔案

```env
DB_USER=
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=
JWT_SECRET=
REDIS_ADDR=localhost:6379
```

### 4️⃣ 執行專案

```bash
go run main.go
```

伺服器預設會在 `localhost:8080` 上運行。

---

## ☸️ Kubernetes 開發部署（使用 KIND）

### 1️⃣ 建立 kind-config.yaml（放在根目錄）

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 30080
        hostPort: 30080
        protocol: TCP
```

### 2️⃣ 建立叢集與部署

```bash
kind create cluster --name dev-cluster --config kind-config.yaml
docker build -t go-app:latest .
kind load docker-image go-app:latest --name dev-cluster
kubectl apply -f k8s/
```

### 3️⃣ 建立 secret.yaml 在k8s的資料夾（請勿直接 commit 真實密碼）

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
stringData:
  DB_PASSWORD: "1234"
  JWT_SECRET: "secrect.yaml"
```

### 4️⃣ 測試 API 與 Swagger

```bash
http://localhost:30080/swagger/index.html
```

---

## 🔐 Secret 管理建議

請勿將含有明文密碼的 `secret.yaml` 提交至 Git。建議如下：

- 保留 `k8s/secret.example.yaml` 作為範本

- 使用者 clone 專案後自行複製並填寫

- `.gitignore` 加入：

  ```
  k8s/secret.yaml
  .env*
  ```

---

## 📘 API 文件

Swagger UI：

```
http://localhost:8080/swagger/index.html
```

---

## 🔐 API 路由總覽

| 方法 | 路由              | 說明                      |
| ---- | ----------------- | ------------------------- |
| POST | `/register`       | 使用者註冊                |
| POST | `/login`          | 使用者登入，發送 JWT      |
| POST | `/refresh`        | 用 Refresh Token 換新 JWT |
| POST | `/logout`         | 清除登入憑證              |
| GET  | `/user/profile`   | 查詢使用者資訊（需登入）  |
| GET  | `/err/test-panic` | 測試全域錯誤攔截          |

---

## 📁 專案結構

```
go-auth-system/
├── config/
├── controllers/
├── middlewares/
├── models/
├── routes/
├── utils/
├── docs/
├── k8s/
│   ├── app-deployment.yaml
│   ├── db-deployment.yaml
│   ├── redis-deployment.yaml
│   ├── configmap.yaml
│   ├── secret.example.yaml
│   └── ...
├── kind-config.yaml
├── main.go
└── go.mod
```

---

## ✨ TODO

- [ ] 加入 Email 驗證功能
- [ ] 加入角色權限管理
- [ ] 撰寫單元測試
- [ ] 支援 Docker / K8s 多環境部署

---

## 📮 聯絡我

Timmy — [taguo1109@gmail.com]
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

---

## 🚀 快速開始

### 1️⃣ 安裝依賴

```bash
go mod tidy
```

### 2️⃣ 初始化 Swagger 文檔

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 3️⃣ 建立資料庫與 .env 環境設定（如果有）
```yaml
-- .env.local 規格
DB_USER=
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=
JWT_SECRET=
REDIS_ADDR=localhost:6379
```

請自行建立 `users` 資料表，並設置 MySQL & Redis 連線資訊（可寫在 `config` 中）

### 4️⃣ 執行專案

```bash
go run main.go
```

伺服器預設會在 `localhost:8080` 上運行。

---

## 📘 API 文件

啟動專案後，打開瀏覽器：

📎 Swagger UI：[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🔐 API 說明

| 方法 | 路由              | 說明                      |
| ---- | ----------------- | ------------------------- |
| POST | `/register`       | 使用者註冊                |
| POST | `/login`          | 使用者登入，發送 JWT      |
| POST | `/refresh`        | 用 Refresh Token 換新 JWT |
| POST | `/logout`         | 清除登入憑證              |
| GET  | `/user/profile`   | 查詢使用者資訊（需登入）  |
| GET  | `/err/test-panic` | 測試全域錯誤攔截          |

---

## 🧪 統一回傳格式（JsonResult）

所有回傳皆使用以下格式：

```json
{
  "status_code": "0000",
  "msg": "Success",
  "msg_detail": "操作成功",
  "data": {}
}
```

錯誤時範例如：

```json
{
  "status_code": "1001",
  "msg": "帳號已存在",
  "msg_detail": "該用戶已存在",
  "data": null
}
```

---

## 🧰 開發者指令

| 指令             | 說明                      |
| ---------------- | ------------------------- |
| `swag init`      | 生成 `docs/` Swagger 文檔 |
| `go run main.go` | 執行伺服器                |
| `go mod tidy`    | 整理依賴                  |

---

## 📁 專案結構

```
go-auth-system/
├── config/         # 資料庫 & Redis 初始化
├── controllers/    # API 控制器（Login, Register, Profile...）
├── middlewares/    # JWT 中介層、全域錯誤處理
├── models/         # 資料模型與 DTO
├── routes/         # 路由註冊
├── utils/          # JWT、回傳格式處理工具
├── docs/           # Swagger 自動生成
├── main.go         # 主程式入口
└── go.mod
```

---

## ✨ TODO

- [ ] 加入 Email 驗證功能
- [ ] 加入角色權限管理
- [ ] 撰寫單元測試
- [ ] 支援 Docker 部署

---

## 📮 聯絡我

Timmy — [taguo1109@gmail.com]

---
# Database Throughput Lab

这是一个用于测试不同数据库吞吐量的实验项目，采用 Go 语言和标准 MVC 架构编写。
支持 MySQL、PostgreSQL、MongoDB，并且可以选择是否启用 Redis 缓存层来对比吞吐量提升。

## 特性

- 🚀 使用 Go + Gin + Gorm + Mongo-Driver 构建
- 📊 封装了独立的 `pkg/benchmark` 吞吐量测试组件
- 💾 支持 MySQL、PostgreSQL、MongoDB 分组测试
- ⚡️ 支持有缓存 (Redis) 和无缓存的读取吞吐量对比测试
- 🐳 提供 Docker Compose 一键拉起所有数据库和缓存环境

## 快速开始

### 1. 启动基础设施

确保你已经安装了 Docker 和 Docker Compose。在项目根目录下运行：

```bash
docker-compose up -d
```

这将会启动 MySQL (3306), PostgreSQL (5432), MongoDB (27017) 和 Redis (6379)。

### 2. 启动服务

使用 Go 启动 HTTP 服务：

```bash
go run cmd/server/main.go
```
服务默认运行在 `http://localhost:8080`。

## API 测试接口

可以使用 Postman 或 `curl` 测试 API。

### 测试写入吞吐量 (Write Test)

**Endpoint:** `POST /api/v1/benchmark/write`

**Payload:**
```json
{
  "db_type": "mysql",    // 可选: "mysql", "postgres", "mongo"
  "concurrency": 50,     // 并发数
  "duration_sec": 10     // 测试持续时间(秒)
}
```

### 测试读取吞吐量 (Read Test)

**Endpoint:** `POST /api/v1/benchmark/read`

**Payload:**
```json
{
  "db_type": "postgres", // 可选: "mysql", "postgres", "mongo"
  "concurrency": 100,
  "duration_sec": 10,
  "use_cache": true      // true: 启用Redis缓存, false: 直接读库
}
```

### 响应示例

```json
{
    "data": {
        "name": "Read Test - postgres (With Cache)",
        "duration": "10.00123s",
        "total_tasks": 125000,
        "successes": 125000,
        "failures": 0,
        "qps": 12498.45,
        "avg_latency": "7.5ms"
    },
    "message": "Read test completed"
}
```

## 目录结构说明

- `cmd/server/main.go`: 程序入口
- `internal/config/`: 数据库连接池和配置
- `internal/controllers/`: API 路由控制层
- `internal/services/`: 业务逻辑层（整合 Benchmark 和 Repo）
- `internal/repository/`: 数据库访问层，处理带缓存和不带缓存的数据获取
- `internal/models/`: 数据模型定义
- `pkg/benchmark/`: 通用的吞吐量测试组件（方便后续自行扩展不同的 SQL 语句测试）

## 推送到 GitHub

你可以执行以下命令将该项目推送到你的 GitHub 仓库：

```bash
git remote add origin <your-github-repo-url>
git branch -M main
git push -u origin main
```

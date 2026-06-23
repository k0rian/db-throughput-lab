# Workspace

本仓库包含一个用于对比不同数据库（以及可选 Redis 缓存）读写吞吐量的 Go 实验项目，便于在本地快速搭建环境并通过 HTTP API 触发压测。

## 子项目

- `db-throughput-lab/`: Database Throughput Lab（Go + Gin + GORM + MongoDB Driver + Redis）
  - 文档：`db-throughput-lab/README.md`
  - 服务入口：`db-throughput-lab/cmd/server/main.go`

## 快速开始

```bash
cd db-throughput-lab
docker-compose up -d
go run cmd/server/main.go
```

服务默认监听：`http://localhost:8080`

## API 示例

写入吞吐量：

```bash
curl -X POST http://localhost:8080/api/v1/benchmark/write \
  -H "Content-Type: application/json" \
  -d '{"db_type":"mysql","concurrency":50,"duration_sec":10}'
```

读取吞吐量（可选缓存对比）：

```bash
curl -X POST http://localhost:8080/api/v1/benchmark/read \
  -H "Content-Type: application/json" \
  -d '{"db_type":"postgres","concurrency":100,"duration_sec":10,"use_cache":true}'
```

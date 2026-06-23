package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/db-throughput-lab/internal/services"
)

type BenchmarkController struct {
	benchmarkService services.BenchmarkService
}

func NewBenchmarkController(service services.BenchmarkService) *BenchmarkController {
	return &BenchmarkController{
		benchmarkService: service,
	}
}

type TestRequest struct {
	DbType      string `json:"db_type" binding:"required"` // mysql, postgres, mongo
	Concurrency int    `json:"concurrency"`
	DurationSec int    `json:"duration_sec"`
	UseCache    bool   `json:"use_cache"`
}

func (c *BenchmarkController) RunWriteTest(ctx *gin.Context) {
	var req TestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Concurrency <= 0 {
		req.Concurrency = 10
	}
	if req.DurationSec <= 0 {
		req.DurationSec = 10
	}

	duration := time.Duration(req.DurationSec) * time.Second

	result, err := c.benchmarkService.RunWriteTest(req.DbType, req.Concurrency, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Write test completed",
		"data":    result,
	})
}

func (c *BenchmarkController) RunReadTest(ctx *gin.Context) {
	var req TestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Concurrency <= 0 {
		req.Concurrency = 10
	}
	if req.DurationSec <= 0 {
		req.DurationSec = 10
	}

	duration := time.Duration(req.DurationSec) * time.Second

	result, err := c.benchmarkService.RunReadTest(req.DbType, req.Concurrency, duration, req.UseCache)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Read test completed",
		"data":    result,
	})
}

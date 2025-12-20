package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type ILogger interface {
	Infof(format string, args ...interface{})
}

// Logger 返回一个 Gin 访问日志中间件
func Logger(logger ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// === 记录请求体（仅对 POST/PUT/PATCH 等有 body 的请求）===
		var reqBody []byte
		if c.Request.Method != "GET" && c.Request.Method != "DELETE" && c.Request.Method != "HEAD" {
			// 限制读取大小，避免大文件导致内存问题（这里限制 1MB，可自行调整）
			reqBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
			// 重新放入 body，让后续 handler 能正常读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// === 捕获响应体 ===
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		cost := time.Since(start)
		clientIP := c.ClientIP()
		respBody := blw.body.String()

		// 基础日志
		logger.Infof("[ACCESS] %s | %13v | %3d | %s %s | IP: %s \n >> Request Body: \n%s\n   << Response Body: %s\n",
			start.Format("2006/01/02 15:04:05"),
			cost,
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			clientIP,
			string(reqBody),
			string(respBody),
		)

	}
}

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 写入缓冲区用于记录
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

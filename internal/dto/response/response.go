package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}

func Message(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, gin.H{
		"success": true,
		"data":    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}

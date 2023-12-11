package service

import "github.com/gin-gonic/gin"

type ApiServiceInterface interface {
	CheckPermission(c *gin.Context) error
	Validate(c *gin.Context) error
	Process(c *gin.Context) (error, int)
	PostProcess(c *gin.Context) error
}

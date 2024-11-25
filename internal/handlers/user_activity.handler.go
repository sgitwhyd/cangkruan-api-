package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
)

type userActHandler struct {
	*gin.Engine
	service service.UserActService
}

func NewUserActHandler(api *gin.Engine, service service.UserActService) *userActHandler {
	return &userActHandler{
		Engine: api,
		service: service,
	}
}

func (h *userActHandler) LikePost(c *gin.Context) {
	ctx := c.Request.Context()

	paramPostID := c.Param("post_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		data := gin.H{
			"error": errors.New("invalid post id").Error(),
		}
		c.JSON(http.StatusUnprocessableEntity, data)
		return
	}

	var body model.CreateUserActivityRequest
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetInt64("userID")

	err = h.service.Create(ctx, body, userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
}

func (h *userActHandler) RegisterRoute() {
		route := h.Group("user_activity")
		route.Use(middlewares.AuthMiddleware())

		route.PUT("/:post_id/like", h.LikePost)
}
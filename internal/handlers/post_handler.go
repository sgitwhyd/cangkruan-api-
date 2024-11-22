package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
)

type handler struct {
	*gin.Engine
	postService service.PostService
}

func NewPostHandler(api *gin.Engine, postService service.PostService) *handler {
	return &handler{
		Engine:        api,
		postService: postService,
	}
}

func (h *handler) Make(c *gin.Context) {
	ctx := c.Request.Context()

	var body model.CreatePostRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Error().Err(err).Msgf("body required")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(int64)
	err = h.postService.Save(ctx, body, userID)
	if err != nil {
		log.Error().Err(err).Msgf("error create post by user_id: %d", userID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": "Post Created",
	})
}	

func (h *handler) RegisterRoute(){
	route := h.Group("posts")
	route.Use(middlewares.AuthMiddleware())

	route.POST("/", h.Make)
}

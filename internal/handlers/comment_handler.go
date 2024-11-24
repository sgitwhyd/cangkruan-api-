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

type commentHandler struct {
	*gin.Engine
	service service.CommentService
	postSvc service.PostService
}

func NewCommentHandler(api *gin.Engine, service service.CommentService, postSvc service.PostService) *commentHandler {
	return &commentHandler{
		Engine: api,
		service: service,
		postSvc: postSvc,
	}
}

func (h commentHandler) Make(c *gin.Context) {
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

	_, err = h.postSvc.FindByID(ctx, postID)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, data)
		return
	}

	var body model.CreateCommentRequest
	err = c.ShouldBind(&body)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	userID := c.MustGet("userID").(int64)

	err = h.service.Save(postID, ctx, body, userID)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, data)
		return
	}
	data := gin.H{
			"data": "created",
		}
	c.JSON(http.StatusCreated, data)
}

func (h *commentHandler) RegisterRoute(){
	route := h.Group("posts")
	route.Use(middlewares.AuthMiddleware())
	route.POST("/:post_id/comments", h.Make)
}
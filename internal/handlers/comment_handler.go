package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/formater"
)

type commentHandler struct {
	*gin.RouterGroup
	service service.CommentService
	postSvc service.PostService
}

func NewCommentHandler(api *gin.RouterGroup, service service.CommentService, postSvc service.PostService) *commentHandler {
	return &commentHandler{
		RouterGroup: api,
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
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create comment", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := c.GetInt64("userID")
	_, err = h.postSvc.FindByID(ctx, userID, postID)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create comment", http.StatusNotFound, "error", data)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var body model.CreateCommentRequest
	err = c.ShouldBind(&body)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create comment", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	err = h.service.Save(postID, ctx, body, userID)
	if err != nil {
			data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create comment", http.StatusInternalServerError, "error", data)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := formater.APIResponse("success create comment", http.StatusCreated, "success", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *commentHandler) RegisterRoute(){
	route := h.Group("posts")
	route.Use(middlewares.AuthMiddleware())
	route.POST("/:post_id/comments", h.Make)
}
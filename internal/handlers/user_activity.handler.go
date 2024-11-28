package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/formater"
)

type userActHandler struct {
	*gin.RouterGroup
	service service.UserActService
	postService service.PostService
}

func NewUserActHandler(api *gin.RouterGroup, service service.UserActService, postService service.PostService) *userActHandler {
	return &userActHandler{
		RouterGroup: api,
		service: service,
		postService: postService,
	}
}

func (h *userActHandler) LikePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64("userID")


	var params model.GetPostIdParam

	err := c.ShouldBindUri(&params)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	post, err := h.postService.FindByID(ctx, userID, int64(params.PostID))
	if err != nil {
			data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed", http.StatusNotFound, "error", data)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var body model.CreateUserActivityRequest
	err = c.ShouldBind(&body)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	msg := "like post"
	if !body.IsLiked {
		msg = "unlike post"
	}


	err = h.service.Create(ctx, body, userID, post.Post.ID)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse(fmt.Sprintf("failed %s", msg), http.StatusInternalServerError, "error", data)
		c.JSON(http.StatusInternalServerError, response)
		return
	}


	response := formater.APIResponse(fmt.Sprintf("success %s", msg), http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *userActHandler) RegisterRoute() {
		route := h.Group("user_activity")
		route.Use(middlewares.AuthMiddleware())

		route.PUT("/:post_id/like", h.LikePost)
}
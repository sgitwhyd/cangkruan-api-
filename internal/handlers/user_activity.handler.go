package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/formater"
)

type userActHandler struct {
	*gin.RouterGroup
	service service.UserActService
}

func NewUserActHandler(api *gin.RouterGroup, service service.UserActService) *userActHandler {
	return &userActHandler{
		RouterGroup: api,
		service: service,
	}
}

func (h *userActHandler) LikePost(c *gin.Context) {
	ctx := c.Request.Context()

	paramPostID := c.Param("post_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
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

	log.Info().Msgf("body %v", body)


	userID := c.GetInt64("userID")

	err = h.service.Create(ctx, body, userID, postID)
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
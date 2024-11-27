package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/middlewares"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/formater"
)

type authHandler struct {
	*gin.RouterGroup
	service service.AuthService
}

func NewAuthHandler(api *gin.RouterGroup, authSvc service.AuthService) *authHandler {
	return &authHandler{
		RouterGroup: api,
		service: authSvc,
	}
}

func (h *authHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.RefreshTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Required Body", http.StatusUnprocessableEntity, "error", error)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	userID := c.GetInt64("userID")
	accessToken, err := h.service.ValidateRefreshToken(ctx, userID, request)
	if err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Error Getting Access Token", http.StatusInternalServerError, "error", error)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := model.RefreshTokenResponse{
		AccessToken: accessToken,
	}
	response := formater.APIResponse("success retrieve access token", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *authHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Required Body", http.StatusUnprocessableEntity, "error", error)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	err := h.service.SignUp(ctx, request)
	if err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Error Sign Up", http.StatusInternalServerError, "error", error)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := formater.APIResponse("Success Sign Up", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *authHandler) SignIn(c *gin.Context){
	ctx := c.Request.Context()

	var body model.SignInRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Required Body", http.StatusUnprocessableEntity, "error", error)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, refToken, err := h.service.SignIn(ctx, body)
	if err != nil {
		error := gin.H{
			"error" : err.Error(),
		}
		response := formater.APIResponse("Error Sign In", http.StatusInternalServerError, "error", error)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := &model.SignInResponse{
		AccessToken: token,
		RefreshToken: refToken,
	}

	response := formater.APIResponse("Success Sign In", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *authHandler) RegisterRoute(){
	route := h.Group("auth")

	route.POST("/signup", h.SignUp)
	route.POST("/signin", h.SignIn)


	h.POST("/auth/refresh", middlewares.AuthRefreshMiddleware(), h.Refresh)
}
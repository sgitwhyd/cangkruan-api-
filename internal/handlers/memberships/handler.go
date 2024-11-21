package memberships

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/sgitwhyd/cangkruan-api/internal/model"
	service "github.com/sgitwhyd/cangkruan-api/internal/service/memberships"
)

type Handler struct {
	*gin.Engine
	membershipSvc service.MembershipService
}

func NewHandler(api *gin.Engine, membershipSvc service.MembershipService) *Handler {
	return &Handler{
		Engine: api,
		membershipSvc: membershipSvc,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Print("request body not fill up")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	err := h.membershipSvc.SignUp(ctx, request)
	if err != nil {
			log.Print("sign up error", err)
		data := gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, data)
		return
	}

	response := gin.H{
		"result": "Created",
	}
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) RegisterRoute(){
	route := h.Group("memberships")

	route.POST("/signup", h.SignUp)
}
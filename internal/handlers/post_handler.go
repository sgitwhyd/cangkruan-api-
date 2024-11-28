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

type handler struct {
	*gin.RouterGroup
	postService service.PostService
}

func NewPostHandler(api *gin.RouterGroup, postService service.PostService) *handler {
	return &handler{
		RouterGroup:        api,
		postService: postService,
	}
}

func (h *handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64("userID")

	pageSizeStr := c.Query("pageSize")
	pageStr := c.Query("page")

	if pageSizeStr == "" || pageStr == "" {
		error := gin.H{
				"error": "need page and page size param",
			}
		response := formater.APIResponse("failed get post detail", http.StatusBadRequest, "error", error)
			c.JSON(http.StatusBadRequest, response)
		return
	}

	limit, err := strconv.Atoi(pageSizeStr)
	if err != nil {
			error := gin.H{
				"error": err.Error(),
			}
			response := formater.APIResponse("failed get post detail", http.StatusBadRequest, "error", error)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	offset, err := strconv.Atoi(pageStr)
		if err != nil {
			error := gin.H{
				"error": err.Error(),
			}
			response := formater.APIResponse("failed get post detail", http.StatusBadRequest, "error", error)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	posts, err := h.postService.FindAll(ctx, limit, offset, userID)
	if err != nil {
			error := gin.H{
				"error": err.Error(),
			}
			response := formater.APIResponse("failed get all post", http.StatusInternalServerError, "error", error)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	response := formater.APIResponse("success get all post", http.StatusOK, "success", posts)
	c.JSON(http.StatusOK, response)
}

func (h *handler) Make(c *gin.Context) {
	ctx := c.Request.Context()

	var body model.CreatePostRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		error := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create post", http.StatusUnprocessableEntity, "error", error)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userID := c.MustGet("userID").(int64)
	err = h.postService.Save(ctx, body, userID)
	if err != nil {
		error := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("failed create post", http.StatusInternalServerError, "error", error)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

		response := formater.APIResponse("success create post", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
}

func (h *handler) Find(c *gin.Context) {
	ctx := c.Request.Context()

	var postId model.GetPostIdParam

	err := c.ShouldBindUri(&postId)
	if err != nil {
		error := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("error getting post detail", http.StatusBadRequest, "error", error)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := c.GetInt64("userID")
	post, err := h.postService.FindByID(ctx,userID, int64(postId.PostID))
	if err != nil {
		error := gin.H{
			"error": err.Error(),
		}
		response := formater.APIResponse("error getting post detail", http.StatusInternalServerError, "error", error)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := formater.APIResponse("success getting post detail", http.StatusOK, "success", post)
	c.JSON(http.StatusOK, response)

}

func (h *handler) RegisterRoute(){
	route := h.Group("posts")
	route.Use(middlewares.AuthMiddleware())

	route.POST("/", h.Make)
	route.GET("/:post_id", h.Find)
	route.GET("/", h.Get)
}

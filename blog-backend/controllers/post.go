package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/task/go_learn_task/blog-backend/database"
	"github.com/task/go_learn_task/blog-backend/models"
	"github.com/task/go_learn_task/blog-backend/utils"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (pc *PostController) CreatePost(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input")
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	if err := database.DB.Preload("User").First(&post, post.ID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Post created successfully", post)
}

func (pc *PostController) GetAllPosts(c *gin.Context) {
	var posts []models.Post

	query := database.DB.Preload("User")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&posts).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch posts", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Posts fetched successfully", posts)
}

func (pc *PostController) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post fetched successfully", post)
}

func (pc *PostController) UpdatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	if post.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own posts", nil)
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input")
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := database.DB.Save(&post).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post", err)
		return
	}

	if err := database.DB.Preload("User").First(&post, post.ID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post updated successfully", post)
}

func (pc *PostController) DeletePost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	if post.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own posts", nil)
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post deleted successfully", nil)
}

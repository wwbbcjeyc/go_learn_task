package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/task/go_learn_task/blog-backend/database"
	"github.com/task/go_learn_task/blog-backend/models"
	"github.com/task/go_learn_task/blog-backend/utils"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	userID := c.GetUint("userID")
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create comment", err)
		return
	}

	if err := database.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch comment", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Comment created successfully", comment)
}

func (cc *CommentController) GetPostComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Post not found", err)
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch comments", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Comments fetched successfully", comments)
}

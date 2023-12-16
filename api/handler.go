package api

import (
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/common"
	"taskmanager/core"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	user core.UserService
	task core.TaskService
}

func NewAPIHandler(ur core.UserService, task core.TaskService) APIHandler {
	return APIHandler{
		user: ur,
		task: task,
	}
}

func (h *APIHandler) CreateUser(ctx *gin.Context) {

	// Request Validation
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	is_user_admin, err := h.user.CheckIfUserIsAdmin(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if is_user_admin == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Only Admin Can Add Users",
		})
		return
	}

	r := core.CreateUserRequest{}
	if err := ctx.BindJSON(&r); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid Json Data",
		})
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	var user core.User
	user.Email = r.Email
	user.Username = r.UserName
	user.Type = r.Type

	if err := h.user.Create(ctx, user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error Processing Request",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}

func (h *APIHandler) DeleteUser(ctx *gin.Context) {

	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	is_user_admin, err := h.user.CheckIfUserIsAdmin(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if is_user_admin == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Only Admin Can Add Users",
		})
		return
	}

	userId, exists := ctx.GetQuery("userId")
	if exists == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "UserId not found",
		})
		return
	}

	if err := h.user.Delete(ctx, userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error Processing Request",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User Deleted",
	})

	return
}

func (h *APIHandler) GetProfileData(ctx *gin.Context) {
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}
	userData, err := h.user.GetUserProfileData(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error Processing Request",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User Data",
		"data":    userData,
	})
	return
}

func (h *APIHandler) CreateTask(ctx *gin.Context) {
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	is_user_admin, err := h.user.CheckIfUserIsAdmin(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if is_user_admin == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Only Admin Can Add/Edit Tasks",
		})
		return
	}

	r := core.CreateTaskRequest{}
	if err := ctx.BindJSON(&r); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	var taskBody core.Task
	taskBody.Priority = r.Priority
	taskBody.Status = r.Status
	taskBody.Title = r.Title
	taskBody.Description = r.Description
	taskBody.UserID = r.UserID
	taskBody.DueDate = r.DueDate

	h.task.CreateTask(ctx, taskBody)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Assigned Task to User",
	})
	return
}

func (h *APIHandler) UpdateTask(ctx *gin.Context) {
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	is_user_admin, err := h.user.CheckIfUserIsAdmin(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if is_user_admin == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Only Admin Can Add/Edit Tasks",
		})
		return
	}

	taskIDStr := ctx.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid task ID",
		})
		return
	}
	var updatedTask core.Task
	if err := ctx.BindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid Json Data",
		})
		return
	}
	err = h.task.UpdateTask(ctx, uint(taskID), updatedTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
	return
}

func (h *APIHandler) MarkTaskStatus(ctx *gin.Context) {
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	userData, err := h.user.GetUserProfileData(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error Processing Request",
		})
		return
	}

	taskIDStr := ctx.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid task ID",
		})
		return
	}

	isUpdated, err := h.task.MarkTaskStatus(ctx, uint(taskID), userData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	if isUpdated {
		message := fmt.Sprintf("User %s task status updated", userData.Username)
		common.PublishToQueue("TASK_UPDATE_QUEUE", message)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully"})
	return
}

func (h *APIHandler) SearchTask(ctx *gin.Context) {
	requestUser := ctx.GetHeader("user_name")
	if requestUser == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Headers Missing",
		})
		return
	}

	is_user_admin, err := h.user.CheckIfUserIsAdmin(ctx, requestUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if is_user_admin == false {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, core.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Only Admin Can Search Tasks",
		})
		return
	}
	searchString := ctx.Param("search")
	priority := ctx.Param("priority")
	statusStr := ctx.Param("status")
	sortType := ctx.Param("sort")
	fmt.Print(searchString, "searchStringsearchStringsearchString")
	taskData := h.task.SearchTasks(ctx, priority, statusStr, sortType, searchString)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User Data",
		"data":    taskData,
	})
	return
}

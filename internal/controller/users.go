package controller

import (
	"errors"
	"net/http"
	"strconv"

	"cruder/internal/model"
	"cruder/internal/service"
	"cruder/pkg/validation"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetByUsername(ctx, username)
	if err != nil {
		ctx.JSON(code(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user *model.User
	if user, err = c.service.GetByID(ctx, id); err != nil {
		ctx.JSON(code(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) PostUser(ctx *gin.Context) {
	user := new(model.User)

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Post(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (c *UserController) PatchUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user *model.User
	if user, err = c.service.GetByID(ctx, id); err != nil {
		ctx.JSON(code(err), gin.H{"error": err.Error()})
		return
	}

	if err = ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Avoid changing the id if it's included in the body and has a different value
	user.ID = id

	if err = c.service.Patch(ctx, user); err != nil {
		ctx.JSON(code(err), gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err = c.service.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

var badRequest validation.InvalidRequest

func code(err error) int {
	switch {
	case errors.Is(err, validation.ErrUserNotFound):
		return http.StatusNotFound
	case errors.As(err, &badRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

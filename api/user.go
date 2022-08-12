package api

import (
	"database/sql"
	"net/http"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username       string `json:"username" binding:"required,min=6,max=10"`
	FullName       string `json:"full_name" binding:"required,alphaunicode"`
	Email          string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams {
		Username: req.Username,
		HashedPassword: "xyz",
		FullName: req.FullName,
		Email: req.Email,
		TotalExpenses: 0,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	Username       string `uri:"username" binding:"required,min=6,max=10"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}


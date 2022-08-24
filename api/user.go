package api

import (
	"database/sql"
	"net/http"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum,min=6,max=10"`
	FullName       string `json:"full_name" binding:"required,alphaunicode"`
	Email          string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	TotalExpenses     int64     `json:"total_expenses"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hasedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams {
		Username: req.Username,
		HashedPassword: hasedPassword,
		FullName: req.FullName,
		Email: req.Email,
		TotalExpenses: 0,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "users_pkey" {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse {
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		TotalExpenses: user.TotalExpenses,
	}

	ctx.JSON(http.StatusOK, rsp)
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

	rsp := createUserResponse {
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		TotalExpenses: user.TotalExpenses,
	}

	ctx.JSON(http.StatusOK, rsp)
}


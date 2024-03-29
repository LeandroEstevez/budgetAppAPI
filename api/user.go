package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/token"
	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=1,max=15"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	TotalExpenses int64  `json:"total_expenses"`
	AccessToken   string `json:"access_token"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		TotalExpenses: user.TotalExpenses,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	// ctx.Header("Access-Control-Allow-Origin", "*")
	// ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	// ctx.Header("Access-Control-Allow-Origin", "*")
	// ctx.Header("Access-Control-Allow-Methods", "*")
	// ctx.Header("Access-Control-Allow-Headers", "*")

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

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hasedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		TotalExpenses:  0,
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

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResponse{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		TotalExpenses: user.TotalExpenses,
		AccessToken:   accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,min=6,max=10"`
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type logInUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=6,max=10"`
	Password string `json:"password" binding:"required,min=6"`
}

type logInUserResponse struct {
	// AccessToken string       `json:"access_token"`
	User userResponse `json:"user"`
}

func (server *Server) logInUser(ctx *gin.Context) {
	var req logInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// rsp := logInUserResponse{
	// 	// AccessToken: accessToken,
	// 	User: newUserResponse(user),
	// }

	rsp := userResponse{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		TotalExpenses: user.TotalExpenses,
		AccessToken:   accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type deleteUserRequest struct {
	Username string `uri:"username" binding:"required,min=6,max=10"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	// ctx.Header("Access-Control-Allow-Origin", "*")
	// ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	// ctx.Header("Access-Control-Allow-Origin", "*")
	// ctx.Header("Access-Control-Allow-Methods", "*")
	// ctx.Header("Access-Control-Allow-Headers", "*")

	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.store.DeleteUserTx(ctx, authPayload.Username)
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

	ctx.JSON(http.StatusOK, "Deletion Completed")
}

type forgotPasswordRequest struct {
	Username string `json:"username" binding:"required,min=6,max=10"`
}

func (server *Server) forgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("Binding json request")

	user, err := server.store.GetEmail(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("Got User from DB", user)

	resetToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("Made the token")

	var firstName = user.FullName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := util.EmailData{
		URL:       "http://localhost:3001" + "/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 15min)",
		ToEmail:   user.Email,
	}

	fmt.Println("Trying to send the email")

	err = util.SendEmail(&emailData, "resetPassword.html")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	ctx.JSON(http.StatusOK, "You will receive a reset email if username exists")
}

type resetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println("got the req")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println("got the payload")

	hasedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("hashed the password")

	arg := db.ResetPasswordParams{
		Username:       authPayload.Username,
		HashedPassword: hasedPassword,
	}

	fmt.Println("got the arg", arg)

	err = server.store.ResetPassword(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("changed the password")

	ctx.JSON(http.StatusOK, "Password data updated successfully")
}

type updateAccountRequest struct {
	OrigUsername string `json:"origusername" binding:"required,alphanum,min=1,max=15"`
	Username     string `json:"username" binding:"required,alphanum,min=1,max=15"`
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountTxParams{
		OrigUsername: req.OrigUsername,
		Username:     req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
	}
	updateUserResult, err := server.store.UpdateAccountTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updateUserResult.User)
}

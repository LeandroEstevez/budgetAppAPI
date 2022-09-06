package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/token"
	"github.com/gin-gonic/gin"
)

const (
	YYYYMMDD = "2006-01-02"
)

type createEntryRequest struct {
	Name    string    `json:"name" binding:"required,alphaunicode"`
	DueDate string `json:"due_date" binding:"required" time_format:"2006-01-02"`
	Amount  int64     `json:"amount" binding:"required,gt=0"`
}

func (server *Server) addEntry(ctx *gin.Context) {
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dueDate, err := time.Parse(YYYYMMDD, req.DueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddEntryTxParams {
		Username: authPayload.Username,
		Name: req.Name,
		DueDate: dueDate,
		Amount: req.Amount,
	}

	entryResult, err := server.store.AddEntryTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entryResult)
}

type deleteEntryRequest struct {
	ID  int32     `uri:"id" binding:"required,gt=0"`
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.DeleteEntryTxParams {
		Username: authPayload.Username,
		ID: req.ID,
	}

	deleteEntryResult, err := server.store.DeleteEntryTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deleteEntryResult)
}

type getEntriesRequest struct {
	Username       string `form:"username" binding:"required,min=6,max=10"`
}

func (server *Server) getEntries(ctx *gin.Context) {
	var req getEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	entries, err := server.store.GetEntries(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
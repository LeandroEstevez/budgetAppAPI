package api

import (
	"net/http"
	"time"

	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/gin-gonic/gin"
)

const (
	YYYYMMDD = "2006-01-02"
)

type createEntryRequest struct {
	Owner   string    `json:"owner" binding:"required,min=6,max=10"`
	Name    string    `json:"name" binding:"required,alphaunicode"`
	DueDate string `json:"due_date" binding:"required" time_format:"2006-01-02"`
	Amount  int64     `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createEntry(ctx *gin.Context) {
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

	arg := db.CreateEntryParams {
		Owner: req.Owner,
		Name: req.Name,
		DueDate: dueDate,
		Amount: req.Amount,
	}

	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type deleteEntryRequest struct {
	ID  int32     `uri:"id" binding:"required,gt=0"`
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteEntry(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}
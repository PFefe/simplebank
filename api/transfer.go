package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/PFefe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	log.Printf(
		"Creating transfer from account %d to account %d of amount %d %s",
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
		req.Currency,
	)

	if !server.validAccount(
		ctx,
		req.FromAccountID,
		req.Currency,
	) {
		return
	}

	if !server.validAccount(
		ctx,
		req.ToAccountID,
		req.Currency,
	) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(
		ctx,
		arg,
	)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		result,
	)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(
		ctx,
		accountId,
	)
	if err != nil {
		if errors.Is(
			err,
			sql.ErrNoRows,
		) {
			ctx.JSON(
				http.StatusNotFound,
				errorResponse(err),
			)
			return false
		}
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf(
			"account [%d] currency mismatch: %s",
			accountId,
			account.Currency,
		)
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)

		return false
	}
	return true
}

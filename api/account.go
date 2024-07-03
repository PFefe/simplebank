package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := server.store.CreateAccount(
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
		account,
	)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(context *gin.Context) {
	var req getAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	account, err := server.store.GetAccount(
		context,
		req.ID,
	)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		account,
	)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listAccount(context *gin.Context) {
	var req listAccountRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	arg := db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccount(
		context,
		arg,
	)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		accounts,
	)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(context *gin.Context) {
	var req deleteAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	err := server.store.DeleteAccount(
		context,
		req.ID,
	)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{"status": "deleted"},
	)

}

type updateAccountURIRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountJSONRequest struct {
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(context *gin.Context) {
	var uriReq updateAccountURIRequest
	if err := context.ShouldBindUri(&uriReq); err != nil {
		context.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	var jsonReq updateAccountJSONRequest
	if err := context.ShouldBindJSON(&jsonReq); err != nil {
		context.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	arg := db.UpdateAccountParams{
		ID:      uriReq.ID,
		Balance: jsonReq.Balance,
	}
	account, err := server.store.UpdateAccount(
		context,
		arg,
	)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		account,
	)
}

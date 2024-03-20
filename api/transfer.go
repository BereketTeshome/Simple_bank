package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/token"
)

type transferRequest struct{
	FromAccountID 	int64 `json:"from_account_id" binding:"required"`
	ToAccountID 	int64 `json:"to_account_id" binding:"required"`
	Amount 	 		int64 `json:"amount" binding:"required"`
	Currency 		string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(c, req.FromAccountID, req.Currency)
	if !valid{
		return
	}

	authPayload := c.MustGet(authorizationHeaderKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("account does not belong to you")
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	_, valid = server.validAccount(c, req.ToAccountID, req.Currency)
	if !valid{
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency,currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}	
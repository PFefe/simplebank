package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Store interface defines all store methods
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore struct implements Store and provides methods to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new SQLStore
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(
		ctx,
		nil,
	)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf(
				"tx error: %v, rb error: %v",
				err,
				rbErr,
			)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(
		ctx,
		func(q *Queries) error {
			var err error

			result.Transfer, err = q.CreateTransfer(
				ctx,
				CreateTransferParams{
					FromAccountID: arg.FromAccountID,
					ToAccountID:   arg.ToAccountID,
					Amount:        arg.Amount,
				},
			)
			if err != nil {
				log.Printf(
					"Failed to create transfer: %v",
					err,
				)
				return err
			}

			result.FromEntry, err = q.CreateEntry(
				ctx,
				CreateEntryParams{
					AccountID: arg.FromAccountID,
					Amount:    -arg.Amount,
				},
			)
			if err != nil {
				log.Printf(
					"Failed to create from entry: %v",
					err,
				)
				return err
			}

			result.ToEntry, err = q.CreateEntry(
				ctx,
				CreateEntryParams{
					AccountID: arg.ToAccountID,
					Amount:    arg.Amount,
				},
			)
			if err != nil {
				log.Printf(
					"Failed to create to entry: %v",
					err,
				)
				return err
			}

			if arg.FromAccountID < arg.ToAccountID {
				result.FromAccount, result.ToAccount, err = addMoney(
					ctx,
					q,
					arg.FromAccountID,
					-arg.Amount,
					arg.ToAccountID,
					arg.Amount,
				)
			} else {
				result.ToAccount, result.FromAccount, err = addMoney(
					ctx,
					q,
					arg.ToAccountID,
					arg.Amount,
					arg.FromAccountID,
					-arg.Amount,
				)
			}

			if err != nil {
				log.Printf(
					"Failed to update accounts: %v",
					err,
				)
			}
			return err
		},
	)

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(
		ctx,
		AddAccountBalanceParams{
			Amount: amount1,
			ID:     accountID1,
		},
	)
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(
		ctx,
		AddAccountBalanceParams{
			Amount: amount2,
			ID:     accountID2,
		},
	)
	if err != nil {
		return
	}
	return
}

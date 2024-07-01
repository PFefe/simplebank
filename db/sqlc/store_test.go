package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				},
			)
			errs <- err
			results <- result
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(
			t,
			err,
		)
		require.NotEmpty(
			t,
			result,
		)

		// Check transfer
		require.Equal(
			t,
			account1.ID,
			result.Transfer.FromAccountID,
		)
		require.Equal(
			t,
			account2.ID,
			result.Transfer.ToAccountID,
		)
		require.Equal(
			t,
			amount,
			result.Transfer.Amount,
		)
		require.NotZero(
			t,
			result.Transfer.ID,
		)
		require.NotZero(
			t,
			result.Transfer.CreatedAt,
		)
		_, err = store.GetTransfer(
			context.Background(),
			result.Transfer.ID,
		)
		require.NoError(
			t,
			err,
		)

		// Check entries
		require.NotZero(
			t,
			result.FromEntry.ID,
		)
		require.NotZero(
			t,
			result.ToEntry.ID,
		)
		require.Equal(
			t,
			account2.ID,
			result.ToEntry.AccountID,
		)
		require.Equal(
			t,
			amount,
			result.ToEntry.Amount,
		)
		require.NotZero(
			t,
			result.ToEntry.Amount,
		)
		require.NotZero(
			t,
			result.ToEntry.CreatedAt,
		)
		_, err = store.GetEntry(
			context.Background(),
			result.ToEntry.ID,
		)
		require.NoError(
			t,
			err,
		)

		require.Equal(
			t,
			account2.ID,
			result.ToEntry.AccountID,
		)
		require.Equal(
			t,
			amount,
			result.ToEntry.Amount,
		)
		require.NotZero(
			t,
			result.ToEntry.CreatedAt,
		)

		_, err = store.GetEntry(
			context.Background(),
			result.ToEntry.ID,
		)
		require.NoError(
			t,
			err,
		)

	}
	// Check the final updated balances
	updatedAccount1, err := store.GetAccount(
		context.Background(),
		account1.ID,
	)
	require.NoError(
		t,
		err,
	)

	updatedAccount2, err := store.GetAccount(
		context.Background(),
		account2.ID,
	)
	require.NoError(
		t,
		err,
	)

	require.Equal(
		t,
		account1.Balance-int64(n)*amount,
		updatedAccount1.Balance,
	)
	require.Equal(
		t,
		account2.Balance+int64(n)*amount,
		updatedAccount2.Balance,
	)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: fromAccountID,
					ToAccountID:   toAccountID,
					Amount:        amount,
				},
			)
			errs <- err
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(
			t,
			err,
		)

	}

	// Check final UUpdated balances
	updateAccount1, err := store.GetAccount(
		context.Background(),
		account1.ID,
	)
	require.NoError(
		t,
		err,
	)
	updateAccount2, err := store.GetAccount(
		context.Background(),
		account2.ID,
	)
	require.NoError(
		t,
		err,
	)
	require.Equal(
		t,
		account1.Balance,
		updateAccount1.Balance,
	)
	require.Equal(
		t,
		account2.Balance,
		updateAccount2.Balance,
	)

}

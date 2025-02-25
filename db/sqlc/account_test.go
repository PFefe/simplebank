package db

import (
	"context"
	"testing"
	"time"

	"github.com/PFefe/simplebank/util"
	"github.com/stretchr/testify/require"
)

// createRandomAccount is a helper function to create a random account
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(
		context.Background(),
		arg,
	)
	require.NoError(
		t,
		err,
	)
	require.NotEmpty(
		t,
		account,
	)

	require.Equal(
		t,
		arg.Owner,
		account.Owner,
	)
	require.Equal(
		t,
		arg.Balance,
		account.Balance,
	)
	require.Equal(
		t,
		arg.Currency,
		account.Currency,
	)
	require.NotZero(
		t,
		account.ID,
	)
	require.NotZero(
		t,
		account.CreatedAt,
	)

	return account
}

// TestCreateAccount tests the creation of an account
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// TestGetAccount tests retrieving an account
func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(
		context.Background(),
		account1.ID,
	)
	require.NoError(
		t,
		err,
	)
	require.NotEmpty(
		t,
		account2,
	)

	require.Equal(
		t,
		account1.ID,
		account2.ID,
	)
	require.Equal(
		t,
		account1.Owner,
		account2.Owner,
	)
	require.Equal(
		t,
		account1.Balance,
		account2.Balance,
	)
	require.Equal(
		t,
		account1.Currency,
		account2.Currency,
	)
	require.WithinDuration(
		t,
		account1.CreatedAt,
		account2.CreatedAt,
		time.Second,
	)
}

// TestUpdateAccount tests updating an account
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	// Update the account balance
	account2, err := testQueries.UpdateAccount(
		context.Background(),
		arg,
	)
	require.NoError(
		t,
		err,
	)
	require.NotEmpty(
		t,
		account2,
	)

	// Fetch the updated account to verify the changes
	account3, err := testQueries.GetAccount(
		context.Background(),
		account1.ID,
	)
	require.NoError(
		t,
		err,
	)
	require.NotEmpty(
		t,
		account3,
	)

	expectedBalance := account1.Balance + arg.Balance

	require.Equal(
		t,
		account1.ID,
		account3.ID,
	)
	require.Equal(
		t,
		account1.Owner,
		account3.Owner,
	)
	require.Equal(
		t,
		expectedBalance,
		account3.Balance,
	)
	require.Equal(
		t,
		account1.Currency,
		account3.Currency,
	)
	require.WithinDuration(
		t,
		account1.CreatedAt,
		account3.CreatedAt,
		time.Second,
	)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(
		context.Background(),
		arg,
	)
	require.NoError(
		t,
		err,
	)
	require.Len(
		t,
		accounts,
		5,
	)

	for _, account := range accounts {
		require.NotEmpty(
			t,
			account,
		)
	}
}

// TestDeleteAccount tests deleting an account
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(
		context.Background(),
		account1.ID,
	)
	require.NoError(
		t,
		err,
	)

	account2, err := testQueries.GetAccount(
		context.Background(),
		account1.ID,
	)
	require.Error(
		t,
		err,
	)
	require.Empty(
		t,
		account2,
	)
}

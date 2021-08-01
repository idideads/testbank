package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/idideads/simplebank/tool"
	"github.com/stretchr/testify/require"
)

func createFakeAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    tool.RandomOwner(),
		Blance:   tool.RandomAmount(),
		Currency: tool.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Blance, account.Blance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	for i := 0; i < 20; i++ {
		createFakeAccount(t)
	}
	// createFakeAccount(t)
}

func TestGetAccount(t *testing.T) {
	refAccount := createFakeAccount(t)

	account, err := testQueries.GetAccount(context.Background(), refAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, refAccount.Owner, account.Owner)
	require.Equal(t, refAccount.Blance, account.Blance)
	require.Equal(t, refAccount.Currency, account.Currency)
	require.Equal(t, refAccount.ID, account.ID)
	// require.Equal(t, fakeAccount.CreatedAt, account.CreatedAt)
	require.WithinDuration(t, refAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	refAccount := createFakeAccount(t)

	arg := UpdateAccountParams{
		ID:     refAccount.ID,
		Blance: tool.RandomAmount(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, refAccount.ID, account.ID)
	require.Equal(t, refAccount.Owner, account.Owner)
	require.Equal(t, refAccount.Currency, account.Currency)
	require.WithinDuration(t, refAccount.CreatedAt, account.CreatedAt, time.Second)
	require.Equal(t, arg.Blance, account.Blance)
}

func TestDeleteAccount(t *testing.T) {
	refAccount := createFakeAccount(t)

	err := testQueries.DeleteAccount(context.Background(), refAccount.ID)
	require.NoError(t, err)

	findAccount, err := testQueries.GetAccount(context.Background(), refAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, findAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

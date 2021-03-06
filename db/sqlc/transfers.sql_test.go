// Code generated by sqlc. DO NOT EDIT.
// source: transfers.sql

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/idideads/simplebank/tool"
	"github.com/stretchr/testify/require"
)

func createFakeTransfer(t *testing.T) Transfer {
	formAccount := createFakeAccount(t)
	toAccount := createFakeAccount(t)

	arg := CreateTransferParams{
		FromAccountID: formAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        tool.RandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	createFakeTransfer(t)
}

func TestQueries_DeleteTransfer(t *testing.T) {
	refTransfer := createFakeTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), refTransfer.ID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), refTransfer.ID)
	require.Error(t, err)
	require.Empty(t, transfer)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestQueries_GetTransfer(t *testing.T) {
	refTransfer := createFakeTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), refTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, refTransfer.ID, transfer.ID)
	require.Equal(t, refTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, refTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, refTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, refTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestQueries_ListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeTransfer(t)
	}

	arg := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

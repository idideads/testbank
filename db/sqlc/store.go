package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

type callBackTxFunc func(q *Queries) error

func (store *Store) execTx(ctx context.Context, fn callBackTxFunc) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("txErr err: %v, rollbackErr err: %v", err, rbErr)
		}
		return err
	}

	err = tx.Commit()
	return err

}

type TransferTxParm struct {
	FromAccountId int64 `json:"form_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParm) (TransferResult, error) {
	var transferResult TransferResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		transferResult.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// transferResult.FromAccount, err = q.AddAccountBlance(ctx, AddAccountBlanceParams{
		// 	ID:     arg.FromAccountId,
		// 	Blance: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// transferResult.ToAccount, err = q.AddAccountBlance(ctx, AddAccountBlanceParams{
		// 	ID:     arg.ToAccountId,
		// 	Blance: arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		if arg.ToAccountId > arg.FromAccountId {
			transferResult.FromAccount, transferResult.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			transferResult.ToAccount, transferResult.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return err
	})

	return transferResult, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBlance(ctx, AddAccountBlanceParams{ID: accountID1, Blance:  amount1})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBlance(ctx, AddAccountBlanceParams{ID: accountID2, Blance:  amount2})
	if err != nil {
		return
	}

	return
}


package helper

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx context.Context, tx pgx.Tx){
	err := recover()
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		panic(rollbackErr)
	} else {
		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			ErrMsg(commitErr, "failed to commit")
		}
	}
}
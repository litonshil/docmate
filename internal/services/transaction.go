package services

import (
	"docmate/internal/model"
	"docmate/utils/consts"
	"fmt"
	"log/slog"
)

func TransactionRollback(
	txc *model.TXClient,
	entity consts.Entity,
	action consts.Action,
) error {
	if err := txc.Rollback(); err != nil {
		slog.Error(
			fmt.Sprintf(
				"error occurred while transaction rollbacked for %v %v",
				entity,
				action,
			),
			"",
			err.Error(),
		)

		return err
	}

	slog.Info("transaction rolled back successfully ...")

	return nil
}

func TransactionCommit(
	txc *model.TXClient,
	entity consts.Entity,
	action consts.Action,
) error {
	if err := txc.Commit(); err != nil {
		slog.Error(
			fmt.Sprintf(
				"error occurred while %v %v transaction commit",
				entity,
				action,
			),
			"error",
			err,
		)

		return err
	}

	slog.Info("transaction successfully committed...")

	return nil
}

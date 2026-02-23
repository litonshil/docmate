package services

import (
	"docmate/client/logger"
	"docmate/internal/model"
	"docmate/utils/consts"
	"fmt"
)

func TransactionRollback(
	txc *model.TXClient,
	loggeruc logger.LogClient,
	entity consts.Entity,
	action consts.Action,
) error {
	if err := txc.Rollback(); err != nil {
		loggeruc.Error(
			fmt.Sprintf(
				"error occurred while transaction rollbacked for %v %v",
				entity,
				action,
			),
			err,
		)
		return err
	}

	loggeruc.Info("transaction rollbacked successfully ...")
	return nil
}

func TransactionCommit(
	txc *model.TXClient,
	loggeruc logger.LogClient,
	entity consts.Entity,
	action consts.Action,
) error {
	if err := txc.Commit(); err != nil {
		loggeruc.Error(
			fmt.Sprintf(
				"error occurred while %v %v transaction commit",
				entity,
				action,
			),
			err,
		)
		return err
	}

	loggeruc.Info("transaction successfully committed...")
	return nil
}

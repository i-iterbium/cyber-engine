package transaction

import (
	"database/sql"
)

// CompleteTx завершает открытую транзакцию в зависимости от наличия паники и ошибки err
func CompleteTx(tx *sql.Tx, err error) {
	if p := recover(); p != nil {
		tx.Rollback()
	} else if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

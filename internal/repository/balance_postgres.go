package repository

import (
	"fmt"
	"time"

	"github.com/UserNaMEeman/shops/app"
	"github.com/jmoiron/sqlx"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r *BalancePostgres) GetBalance(guid string, totalAccrual float64) (app.Balance, error) {
	var withdrawn float64
	balnce := app.Balance{}
	queryOrder := fmt.Sprintf("SELECT withdrawn FROM %s WHERE user_guid = $1", balanceTable)
	row := r.db.QueryRow(queryOrder, guid) //(queryOrder, guid)
	if err := row.Scan(&withdrawn); err != nil {
		return app.Balance{}, err
	}
	balnce.Current = totalAccrual - withdrawn
	balnce.Withdrawn = withdrawn
	fmt.Println("postgre balance: ", balnce)
	return balnce, nil
	// a := app.Balance{totalAccrual, 0}
	// return a, nil
}

func (r *BalancePostgres) UsePoints(guid string, buy app.Buy) error {
	now := time.Now()
	timeNow := now.Format(time.RFC3339)
	query := fmt.Sprintf("INSERT INTO %s (order_buy, sum, date_buy) values ($1, $2, $3)", buysTable)
	_, err := r.db.Exec(query, buy.Order, buy.Sum, timeNow) //.QueryRow(query, userGUID, orderNumber)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("UPDATE %s SET withdrawn = withdrawn + $1 WHERE user_guid = $2", balanceTable)
	_, err = r.db.Exec(query, buy.Sum, guid)
	if err != nil {
		return err
	}
	return nil
}

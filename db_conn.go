package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func insertPoint(tx Transaction, ctx context.Context) (*PointDB, error) {
	res, err := tx.ExecContext(ctx, "INSERT INTO points (point) VALUES (0);")
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 1 {
		return &PointDB{ID: int32(lastID), Point: 0}, nil
	} else {
		return nil, fmt.Errorf("affected is not 0")
	}
}

func getPoint(tx Transaction, ctx context.Context) (*PointDB, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id, point FROM points FOR UPDATE;")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	result := make([]*PointDB, 0)
	for rows.Next() {
		t := &PointDB{}
		err := rows.Scan(&t.ID, &t.Point)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	if len(result) == 1 {
		return result[0], nil
	} else {
		return nil, nil
	}
}

func updatePoint(tx Transaction, ctx context.Context, ID int32, point int64) (*PointDB, error) {
	query := `
			UPDATE points 
			SET point = ? 
			WHERE id = ?;`

	res, err := tx.ExecContext(
		ctx, query, point, ID)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected > 1 {
		return nil, fmt.Errorf("more than 1 records affected")
	}

	p, err := getPoint(tx, ctx)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getTransactions(tx Transaction, ctx context.Context, limit int) ([]*TransactionDB, error) {
	var query string
	var rows *sql.Rows
	var err error
	if limit == 0 {
		query = `SELECT ` + "`id`, `date`, `previous`, `change`, `final`" +
			`FROM transactions
				ORDER BY date ASC;`

		rows, err = tx.QueryContext(ctx, query)
	} else {
		query = `SELECT ` + "`id`, `date`, `previous`, `change`, `final`" +
			`FROM (
					SELECT * 
					FROM transactions
					ORDER BY date DESC LIMIT ? OFFSET 0
				) sub ORDER BY date ASC;`

		rows, err = tx.QueryContext(ctx, query, limit)
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	result := make([]*TransactionDB, 0)
	for rows.Next() {
		t := &TransactionDB{}
		err := rows.Scan(&t.ID, &t.DateTime, &t.Previous, &t.Change, &t.Final)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func insertTransaction(tx Transaction, ctx context.Context, previous int64, change int64, final int64) (*TransactionDB, error) {
	query := `INSERT INTO transactions (` + "`date`, `previous`, `change`, `final`" + `) VALUES (?, ?, ?, ?);`

	currentTime := time.Now()

	res, err := tx.ExecContext(ctx, query, currentTime, previous, change, final)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &TransactionDB{
		ID:       int32(lastID),
		DateTime: currentTime,
		Previous: &previous,
		Change:   change,
		Final:    final,
	}, nil
}

func insertAltTransaction(tx Transaction, ctx context.Context, previous int64, change int64, final int64) (*TransactionDB, error) {
	query := `
INSERT INTO transactions (` + "`previous_date`, `date`, `previous`, `change`, `final`" + `)
(SELECT ` + "`date`, ?, `final`, ?, (`final` + ?) FROM transactions ORDER BY `date` DESC LIMIT 1)" +
		`UNION
(SELECT NULL as ` + "`previous_date`, ? as `date`, NULL as `previous`, ? as `change`, ? as `final`)" +
		`LIMIT 1;
`

	currentTime := time.Now()

	res, err := tx.ExecContext(ctx, query, currentTime, change, change, currentTime, change, change)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &TransactionDB{
		ID:       int32(lastID),
		DateTime: currentTime,
		Previous: &previous,
		Change:   change,
		Final:    final,
	}, nil
}

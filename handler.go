package main

import (
	"github.com/cenkalti/backoff"
	"github.com/labstack/echo/v4"
	"net/http"
)

func handleGET(ctx echo.Context) error {
	var point *PointDB
	var transactions []*TransactionDB
	var err error

	err = scope.CreateNew(func(tx Transaction) error {
		var err error
		point, err = getPoint(tx, ctx.Request().Context())

		if err != nil {
			return err
		}

		if point == nil {
			point, err = insertPoint(tx, ctx.Request().Context())
		}

		transactions, err = getTransactions(tx, ctx.Request().Context(), 0)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusInternalServerError, &GetResponse{OK: false, Message: &message})
	}

	return ctx.JSON(http.StatusOK, &GetResponse{OK: true, Value: &point.Point, Transactions: transactions})
}

func handlePOST(ctx echo.Context) error {
	var request Request
	err := ctx.Bind(&request)
	if err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusInternalServerError, &PostResponse{
			OK:      false,
			Message: &message,
		})
	}

	var trx *TransactionDB

	//err = scope.CreateNew(func(tx Transaction) error {
	//	operation := func() error {
	//		var err error
	//		point, err := getPoint(tx, ctx.Request().Context())
	//
	//		if err != nil {
	//			return err
	//		}
	//
	//		if point == nil {
	//			point, err = insertPoint(tx, ctx.Request().Context())
	//		}
	//
	//		point, err = updatePoint(tx, ctx.Request().Context(), point.ID, point.Point+int64(request.Value))
	//		if err != nil {
	//			return err
	//		}
	//
	//		trx, err = insertAltTransaction(tx, ctx.Request().Context(), 0, int64(request.Value), 0+int64(request.Value))
	//		if err != nil {
	//			return err
	//		}
	//
	//		return nil
	//	}
	//
	//	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	//	if err != nil {
	//		return err
	//	}
	//
	//	return nil
	//})

	operation := func() error {
		err = scope.CreateNew(func(tx Transaction) error {
			var err error
			point, err := getPoint(tx, ctx.Request().Context())

			if err != nil {
				return err
			}

			if point == nil {
				point, err = insertPoint(tx, ctx.Request().Context())
			}

			point, err = updatePoint(tx, ctx.Request().Context(), point.ID, point.Point+int64(request.Value))
			if err != nil {
				return err
			}

			trx, err = insertAltTransaction(tx, ctx.Request().Context(), 0, int64(request.Value), 0+int64(request.Value))
			if err != nil {
				return err
			}

			return nil
		})

		return err
	}

	err = backoff.Retry(operation, backoff.NewExponentialBackOff())

	if err != nil {
		message := err.Error()
		return ctx.JSON(http.StatusInternalServerError, &PostResponse{
			OK:      false,
			Message: &message,
		})
	}

	return ctx.JSON(http.StatusOK, &PostResponse{OK: true, Transaction: trx})
}

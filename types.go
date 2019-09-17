package main

import "time"

type PointDB struct {
	ID    int32 `json:"id"`
	Point int64 `json:"point"`
}

type TransactionDB struct {
	ID       int32     `json:"id"`
	DateTime time.Time `json:"date_time"`
	Previous int64     `json:"previous"`
	Change   int64     `json:"change"`
	Final    int64     `json:"final"`
}

type Request struct {
	Value int32 `json:"value"`
}

type PostResponse struct {
	OK          bool           `json:"ok"`
	Message     *string        `json:"message"`
	Transaction *TransactionDB `json:"transaction"`
}

type GetResponse struct {
	OK           bool             `json:"ok"`
	Message      *string          `json:"message"`
	Value        *int64           `json:"value"`
	Transactions []*TransactionDB `json:"transactions"`
}

package domain

type TransactionFromFront struct {
	To     int64   `json:"to"`
	Amount float64 `json:"amount"`
}

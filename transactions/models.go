package transactions

import "time"

type TransactionDto struct {
	ID          string
	Date        time.Time
	Description string
	Amount      float64
}

// "Transaction Date","Transaction Cleared Date","Transaction Type","Transaction Description","Transaction Amount"
type lloydsTransaction struct {
	Date        time.Time  `csv:"Transaction Date"`
	ClearedDate *time.Time `csv:"Transaction Cleared Date"`
	Type        string     `csv:"Transaction Type"`
	Description string     `csv:"Transaction Description"`
	Amount      float64    `csv:"Transaction Amount"`
}

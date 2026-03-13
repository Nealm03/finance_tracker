package transactions

import (
	"fmt"
	"math/big"
	"time"
)

type TransactionDto struct {
	ID          string
	Date        time.Time
	Description string
	AmountPence big.Int
}

type lloydsTransaction struct {
	Date        transactionDate  `csv:"Transaction Date"`
	ClearedDate *transactionDate `csv:"Transaction Cleared Date"`
	Type        string           `csv:"Transaction Type"`
	Description string           `csv:"Transaction Description"`
	Amount      float64          `csv:"Transaction Amount"`
}

type transactionDate time.Time

func (t *transactionDate) UnmarshalCSV(data []byte) error {
	parsedDate, err := time.Parse("02/01/2006", string(data))
	if err != nil {
		return fmt.Errorf("failed to parse transaction date:: %w", err)

	}

	*t = transactionDate(parsedDate)
	return nil
}

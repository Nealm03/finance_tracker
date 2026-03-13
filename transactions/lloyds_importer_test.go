package transactions_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Nealm03/finance_tracker/transactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LloydsImporterSuite struct {
	suite.Suite
}

func (suite *LloydsImporterSuite) createImporter(fs fstest.MapFS, name string) *transactions.LloydsImporter {
	imp, err := transactions.NewLLloydsImporter(name, fs, false)
	suite.Require().NoError(err)
	return imp
}

func (suite *LloydsImporterSuite) TestNewImporterNotNil() {
	fs := fstest.MapFS{"foo.csv": {Data: []byte("a,b")}}
	imp := suite.createImporter(fs, "foo.csv")
	assert.NotNil(suite.T(), imp)
}

func (suite *LloydsImporterSuite) TestImportReturnsExpectedTransactions() {
	// header fields correspond to struct tags
	header := `"Transaction Date","Transaction Cleared Date","Transaction Type","Transaction Description","Transaction Amount"`
	rawCsvRows := []string{
		`"05/02/2026","06/02/2026","Contactless purchase","MARKS&SPENCER PLC SACA","24.85"`,
		`"05/02/2026","06/02/2026","Contactless purchase","MORRISONS DAILY","7.35"`,
	}
	expected := []transactions.TransactionDto{
		suite.rawDataToDto("1", rawCsvRows[0]),
		suite.rawDataToDto("2", rawCsvRows[1]),
	}

	csvData := fmt.Sprintf(
		"%s\n%s",
		header,
		strings.Join(rawCsvRows, "\n"),
	)

	fs := fstest.MapFS{"ledger.csv": {Data: []byte(csvData)}}
	imp := suite.createImporter(fs, "ledger.csv")

	got, err := imp.Import(context.Background(), "ledger.csv")

	suite.NoError(err)

	suite.ElementsMatch(expected, got)
}

func (suite *LloydsImporterSuite) rawDataToDto(id string, rawDataLine string) transactions.TransactionDto {

	removeEscapedQuotes := func(rawVal string) string {
		return strings.ReplaceAll(rawVal, `"`, "")
	}
	vals := strings.Split(rawDataLine, ",")
	suite.GreaterOrEqual(vals, 1, "expected raw mock data be comma delimited")
	rawCreatedDate := removeEscapedQuotes(vals[0])

	parsedDate, err := time.Parse("02/01/2006", rawCreatedDate)
	suite.NoError(err, "expected mock data first col to be valid ISO date")

	parsedAmount, err := strconv.ParseFloat(
		removeEscapedQuotes(vals[4]),
		64,
	)
	suite.NoError(err, "expected mock data 5th col to be valid float")
	return transactions.TransactionDto{
		Date:        parsedDate,
		ID:          id,
		Description: removeEscapedQuotes(vals[3]),
		Amount:      parsedAmount,
	}
}

func TestLloydsImporterSuite(t *testing.T) {
	suite.Run(t, new(LloydsImporterSuite))
}

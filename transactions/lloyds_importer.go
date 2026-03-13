package transactions

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

type LloydsImporter struct {
	fileHandle fs.File
}

func (importer *LloydsImporter) Import(ctx context.Context, filePath string) ([]TransactionDto, error) {
	br := bufio.NewReader(importer.fileHandle)
	r, _, err := br.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if r != '\uFEFF' {
		br.UnreadRune()
	}

	data, err := io.ReadAll(br)
	if err != nil {
		return nil, fmt.Errorf("failed to read file:: %w", err)
	}

	// unmarshal into intermediate struct
	var raw []lloydsTransaction
	if err := csvutil.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	// convert to DTOs
	out := make([]TransactionDto, 0, len(raw))
	for index, r := range raw {
		dto := TransactionDto{
			Date:        time.Time(r.Date),
			Description: r.Description,
			AmountPence: *big.NewInt(
				int64(r.Amount) * 100,
			),
		}
		id, err := GenerateIdHash(dto)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hash, row: %d, err:: %w", index, err)
		}
		dto.ID = id
		out = append(out, dto)
	}

	return out, nil
}

func GenerateIdHash(transaction TransactionDto) (string, error) {
	hasher := sha256.New()
	sections := []string{
		fmt.Sprintf("|%d|%s", len(transaction.Date.String()), transaction.Date.String()),
		fmt.Sprintf("|%d|%s", len(transaction.AmountPence.String()), transaction.AmountPence.String()),
		fmt.Sprintf("|%d|%s", len(transaction.Description), transaction.Description),
	}

	_, err := hasher.Write(
		[]byte(strings.Join(sections, "^")),
	)

	if err != nil {
		return "", fmt.Errorf("failed to hash transaction")
	}

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

func NewLLloydsImporter(fileName string, fileSys fs.FS, abortOnMalformedLine bool) (*LloydsImporter, error) {
	fileInfo, err := fs.Stat(fileSys, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to find file: %w", err)
	}

	if fileInfo.IsDir() {
		return nil, errors.New("file path provided is a directory, are you sure you provided the correct path?")
	}

	if !strings.HasSuffix(fileInfo.Name(), ".csv") {
		return nil, errors.New("file is not the expected type of csv")
	}
	file, err := fileSys.Open(fileName)

	if err != nil {
		return nil, fmt.Errorf("failed to open file:: %w", err)
	}

	return &LloydsImporter{
		fileHandle: file,
	}, nil
}

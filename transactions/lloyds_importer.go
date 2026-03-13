package transactions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/jszwec/csvutil"
)

type LloydsImporter struct {
	fileHandle fs.File
}

func (importer *LloydsImporter) Import(ctx context.Context, filePath string) ([]TransactionDto, error) {
	data, err := io.ReadAll(importer.fileHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// unmarshal into intermediate struct
	var raw []lloydsTransaction
	if err := csvutil.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	// convert to DTOs
	out := make([]TransactionDto, 0, len(raw))
	for i, r := range raw {
		out = append(out, TransactionDto{
			ID:          fmt.Sprintf("%d", i+1),
			Date:        r.Date,
			Description: r.Description,
			Amount:      r.Amount,
		})
	}

	return out, nil
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

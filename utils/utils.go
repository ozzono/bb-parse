package utils

import (
	"bb-parse/internal/models"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func TrimSpaces(input string) string {
	for strings.Contains(input, "  ") {
		input = strings.ReplaceAll(input, "  ", " ")
	}
	input = strings.TrimPrefix(input, " ")
	input = strings.TrimSuffix(input, " ")
	return input
}

func CSVWriter(input []*models.Record, outputFile string) error {
	records := [][]string{
		{"date", "description", "value"},
	}
	for i := range input {
		records = append(records,
			[]string{input[i].Date, input[i].Description, input[i].Category, fmt.Sprintf("%.2f", input[i].Value)},
		)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrap(err, "os.Create output file")
	}
	defer f.Close()

	w := csv.NewWriter(f)

	err = w.WriteAll(records)
	if err != nil {
		return errors.Wrap(err, "w.WriteAll records")
	}

	return nil
}

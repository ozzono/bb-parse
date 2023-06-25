package main

import (
	"bb-parse/internal/db"
	"bb-parse/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRow(t *testing.T) {
	sample := `06.04.2023J B K  GESTAO DE ESTAC SAO PAULO     BR               20,00        0,00`
	expected := models.Record{
		Date:        "2023-04-06",
		Description: "j b k gestao de estac sao paulo",
		Value:       20,
	}

	c, err := db.NewDB()
	assert.NoError(t, err)

	row, err := c.ParseAndCompare(sample)
	assert.NoError(t, err)
	assert.Equal(t, expected.Date, row.Date)
	assert.Equal(t, expected.Description, row.Description)
	assert.Equal(t, expected.Category, row.Category)
	assert.Equal(t, expected.Value, row.Value)
}

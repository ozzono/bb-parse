package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRow(t *testing.T) {
	sample := `06.04.2023J B K  GESTAO DE ESTAC SAO PAULO     BR               20,00        0,00`
	expected := row{
		date:        "2023-04-06",
		description: "J B K GESTAO DE ESTAC SAO PAULO",
		value:       20,
	}
	row := parseRow(sample)
	assert.Equal(t, expected.date, row.date)
	assert.Equal(t, expected.description, row.description)
	assert.Equal(t, expected.value, row.value)
}

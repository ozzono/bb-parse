package main

import (
	"bb-parse/utils"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type row struct {
	date        string // format YYYY-MM-DD
	description string
	value       float64
}

var (
	samples = []string{
		"16.05.2023PG JUSBRASIL            557130352528 BR               84,00        0,00",
		"04.04.2023MercPag TEARMANUAL     OSASCO        BR              212,80        0,00",
		"17.04.2023MEDMANIPULADBUENOS AY   551199390353 BR              367,80        0,00",
		"13.04.2023IKESAKI Ikesaki Cosme  SAO PAULO     BR              -27,69        0,00",
	}
	datePattern = `\d{2}\.\d{2}\.\d{4}`
	file        string
	write       bool
	outputFile  string
)

func init() {
	flag.StringVar(&file, "f", "", "sets the path to the text file to parsed")
	flag.BoolVar(&write, "w", false, "sets the path to the text file to parsed")
	flag.StringVar(&outputFile, "o", "samples/out.csv", "sets the path to the text file to parsed")
}

func validate() {
	flag.Parse()
	if file == "" {
		log.Println("invalid file path - can't be nil")
		return
	}
	log.Println("file path --------", file)
	log.Println("outputFile path --", outputFile)
}

func main() {
	validate()
	f, err := readFile()
	if err != nil {
		log.Println(errors.Wrap(err, "readFile"))
		return
	}
	rows := parseFile(f)
	csvWriter(rows)
}

func readFile() (string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", errors.Wrap(err, "ioutil.ReadFile -f file")
	}
	return string(f), nil
}

func parseFile(input string) []row {
	rows := []row{}
	i := 0
	for _, v := range strings.Split(input, "\n") {
		// fmt.Println(v[49:69])
		v = strings.TrimSuffix(v, "\r")
		if len(v) != 81 {
			continue
		}
		match, err := regexp.Match(datePattern, []byte(v[:10]))
		if err != nil {
			continue
		}
		if match {
			rows = append(rows, parseRow(v))
			i++
			fmt.Println(v)
		}

		if i == 10 {
			break
		}
	}
	return rows
}

func parseRow(input string) row {
	out := row{}
	splittedDate := strings.Split(input[:10], ".")
	out.date = strings.Join([]string{splittedDate[2], splittedDate[1], splittedDate[0]}, "-")
	out.description = utils.TrimSpaces(input[10:47])
	v, _ := strconv.ParseFloat(strings.ReplaceAll(utils.TrimSpaces(input[49:69]), ",", "."), 64)
	out.value = v
	return out
}

func csvWriter(input []row) error {
	records := [][]string{
		{"date", "description", "value"},
	}
	for i := range input {
		records = append(records,
			[]string{input[i].date, input[i].description, fmt.Sprintf("%.2f", input[i].value)},
		)
	}
	f, err := os.Create("out.csv")
	if err != nil {
		return errors.Wrap(err, "os.Create output file")
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, r := range records {
		if err := w.Write(r); err != nil {
			return errors.Wrap(err, "w.Write")
		}
	}
	return nil
}

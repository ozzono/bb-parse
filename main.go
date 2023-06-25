package main

import (
	"bb-parse/internal/db"
	"bb-parse/internal/models"
	"bb-parse/utils"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type row struct {
	date        string // format YYYY-MM-DD
	description string
	value       float64
}

var (
	datePattern = `\d{2}\.\d{2}\.\d{4}`
	file        string
	write       bool
	outputFile  string
)

func init() {
	godotenv.Load(".env")
	flag.StringVar(&file, "f", "", "sets the path to the text file to parsed")
	flag.BoolVar(&write, "w", false, "sets the path to the text file to parsed")
	flag.StringVar(&outputFile, "o", "samples/out.csv", "sets the path to the output csv file")
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
	rows, err := parseFile(f)
	if err != nil {
		log.Println(errors.Wrap(err, "parseFile"))
		return
	}
	err = utils.CSVWriter(rows, outputFile)
	if err != nil {
		log.Println(errors.Wrap(err, "utils.CSVWriter"))
	}
}

func readFile() (string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return "", errors.Wrap(err, "os.ReadFile -f file")
	}
	return string(f), nil
}

func parseFile(input string) ([]*models.Record, error) {
	rows := []*models.Record{}
	c, err := db.NewDB()
	if err != nil {
		return nil, errors.Wrap(err, "db.NewDB")
	}
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
			r, err := c.ParseAndCompare(v)
			if err != nil {
				return nil, errors.Wrap(err, "c.ParseAndCompare")
			}
			rows = append(rows, r)
		}
	}
	return rows, nil
}

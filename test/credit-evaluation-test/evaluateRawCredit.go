package main

import (
	credit_evaluation "credit-evaluation/application-gateway/credit-evaluation"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func processCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	outputFile, err := os.Create("./test/credit-evaluation-test/credit_evaluation_test_data.csv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers := append(records[0],
		"Raw Evaluation Score",
		"Raw Evaluation Result",
		"Evaluation Time (ns)")
	if err := writer.Write(headers); err != nil {
		return err
	}

	for i, record := range records[1:] {
		creditScore, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}

		dti, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return err
		}

		start := time.Now()
		score := credit_evaluation.CreditEvaluation(19, 100*1000*1000+1, float64(creditScore), dti)
		elapsed := time.Since(start).Nanoseconds()

		result := credit_evaluation.Sigmoid(score)

		newRecord := append(
			record,
			strconv.FormatFloat(score, 'f', 2, 64),
			strconv.FormatFloat(result, 'f', 2, 64),
			strconv.FormatInt(elapsed, 10),
		)

		if err := writer.Write(newRecord); err != nil {
			return err
		}

		fmt.Printf("Processed record %d: Credit Score = %d, DTI = %.2f, Result = %s\n", i+1, creditScore, dti, result)
	}

	return nil
}

func maidn() {
	filename := "./test/credit-evaluation-test/credit_evaluation_test_data.csv"
	if err := processCSV(filename); err != nil {
		panic(err)
	}

	println("Results saved to credit_evaluation_results.csv")
}

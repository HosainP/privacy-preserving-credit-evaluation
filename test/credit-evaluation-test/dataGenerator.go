package main

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateTestData(numSamples int) [][]string {
	rand.Seed(time.Now().UnixNano())
	var data [][]string

	for i := 0; i < numSamples; i++ {
		creditScore := rand.Intn(551) + 300
		dti := rand.Float64()*0.5 + 0.1
		dtiStr := strconv.FormatFloat(dti, 'f', 2, 64)

		data = append(data, []string{strconv.Itoa(creditScore), dtiStr})
	}

	return data
}

func saveToCSV(data [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Credit Score", "DTI"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func main2() {
	testData := generateTestData(500)

	filename := "./test/credit-evaluation-test/credit_evaluation_test_data.csv"
	if err := saveToCSV(testData, filename); err != nil {
		panic(err)
	}

	println("Data saved to", filename)
}

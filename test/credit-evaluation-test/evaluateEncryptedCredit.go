package main

import (
	credit_evaluation "credit-evaluation/application-gateway/credit-evaluation"
	"credit-evaluation/application-gateway/encryption"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func processCSV2(filename string, helper *encryption.CKKSHelper) error {
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
		"Encrypted Evaluation Score",
		"Encrypted Evaluation Result",
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

		// Encrypt credit score and DTI
		creditScoreEnc := helper.EncryptPu(float64(creditScore))
		dtiEnc := helper.EncryptPu(dti)

		start := time.Now()
		// Perform homomorphic credit evaluation
		resultEnc, _ := helper.CreditEvaluation(helper, nil, nil, creditScoreEnc, dtiEnc)

		// Decrypt the result
		score := helper.Decrypt(resultEnc)
		elapsed := time.Since(start).Nanoseconds()

		result := credit_evaluation.Sigmoid(score)

		newRecord := append(record,
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

func maing() {
	helper := encryption.NewCKKSHelper()
	filename := "./test/credit-evaluation-test/credit_evaluation_test_data.csv"
	if err := processCSV2(filename, helper); err != nil {
		panic(err)
	}

	println("Results saved to credit_evaluation_results.csv")
}

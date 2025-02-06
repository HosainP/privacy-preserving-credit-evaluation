package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	filename := "./test/credit-evaluation-test/credit_evaluation_test_data.csv"

	original, generated, err := readColumns(filename, "Raw Evaluation Score", "Encrypted Evaluation Score")
	if err != nil {
		panic(err)
	}

	mae := meanAbsoluteError(original, generated)
	mse := meanSquaredError(original, generated)

	fmt.Println("Evaluation Score Report:")
	fmt.Printf("Mean Absolute Error (MAE): %.10f\n", mae)
	fmt.Printf("Mean Squared Error (MSE): %.10f\n", mse)

	///////

	//original, generated, err = readColumns(filename, "Raw Evaluation Score", "Encrypted Evaluation Score")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for i := 0; i < len(original); i++ {
	//	if original[i] != generated[i] {
	//		fmt.Printf("Original: %d - Generated: %d\n", original[i], generated[i])
	//	}
	//}
	//
	//fmt.Println("------------------")
	//
	original, generated, err = readColumns(filename, "Raw Evaluation Result", "Encrypted Evaluation Result")
	if err != nil {
		panic(err)
	}
	//
	//for i := 0; i < len(original); i++ {
	//	if original[i] != generated[i] {
	//		fmt.Printf("Original: %d - Generated: %d\n", original[i], generated[i])
	//	}
	//}

	mae = meanAbsoluteError(original, generated)
	mse = meanSquaredError(original, generated)

	fmt.Println("Evaluation Result Report:")
	fmt.Printf("Mean Absolute Error (MAE): %.10f\n", mae)
	fmt.Printf("Mean Squared Error (MSE): %.10f\n", mse)

	///////

	original, generated, err = readColumns(filename, "Raw Evaluation Time (ns)", "Encrypted Evaluation Time (ns)")
	if err != nil {
		panic(err)
	}

	ratio, _ := averageRatio(original, generated)

	sum := 0.0
	for i := 0; i < len(original); i++ {
		sum = sum + original[i]
	}
	println(sum / 200)

	sum = 0.0
	for i := 0; i < len(generated); i++ {
		sum = sum + generated[i]
	}
	println(sum / 200)

	fmt.Println("Average Time Ratio Report")
	fmt.Printf("Average Time Ratio Report: %.4f\n", ratio)

}

// Read two columns by name from a CSV file
func readColumns(filename, col1, col2 string) ([]float64, []float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("the file must have at least one row of data")
	}

	// Find the indices of the desired columns
	header := records[0]
	col1Idx, col2Idx := -1, -1
	for i, name := range header {
		if name == col1 {
			col1Idx = i
		} else if name == col2 {
			col2Idx = i
		}
	}

	if col1Idx == -1 || col2Idx == -1 {
		return nil, nil, fmt.Errorf("columns %q and/or %q not found in the header", col1, col2)
	}

	// Extract data for the specified columns
	var col1Values, col2Values []float64
	for _, record := range records[1:] {
		val1, err := strconv.ParseFloat(record[col1Idx], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing value in column %q: %v", col1, err)
		}
		val2, err := strconv.ParseFloat(record[col2Idx], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing value in column %q: %v", col2, err)
		}
		col1Values = append(col1Values, val1)
		col2Values = append(col2Values, val2)
	}

	return col1Values, col2Values, nil
}

// Calculate Mean Absolute Error (MAE)
func meanAbsoluteError(original, generated []float64) float64 {
	sum := 0.0
	for i := range original {
		sum += math.Abs(original[i] - generated[i])
	}
	return sum / float64(len(original))
}

// Calculate Mean Squared Error (MSE)
func meanSquaredError(original, generated []float64) float64 {
	sum := 0.0
	for i := range original {
		diff := original[i] - generated[i]
		sum += diff * diff
	}
	return sum / float64(len(original))
}

// Calculate the average ratio of y to x
func averageRatio(x, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0, fmt.Errorf("x and y must have the same length")
	}

	sum := 0.0
	for i := range x {
		if x[i] == 0 {
			return 0, fmt.Errorf("division by zero at index %d", i)
		}
		sum += y[i] / x[i]
	}

	return sum / float64(len(x)), nil
}

// Calculate Standard Deviation (SD)
func standardDeviation(original, generated []float64) float64 {
	if len(original) == 0 {
		return 0
	}

	sum := 0.0
	for i := range original {
		diff := original[i] - generated[i]
		sum += diff * diff
	}

	variance := sum / float64(len(original))
	return math.Sqrt(variance)
}

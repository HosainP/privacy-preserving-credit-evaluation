package credit_evaluation

import (
	"math"
)

const MinSalary = 10 * 1000 * 1000
const MinAge = 18

const MaxCreditScore = 850
const MinCreditScore = 300
const MinDTI = 0.01

const W0 = 0.5
const W1 = 0.5

func CreditEvaluation(age int, salary int, creditScore float64, dti float64) float64 {
	if !satisfyPreselection(age, salary) {
		return -1
	}
	return calcScore(creditScore, dti)
}

func satisfyPreselection(age int, salary int) bool {
	return salary > MinSalary && age > MinAge
}

func calcScore(creditScore float64, dti float64) float64 {
	// Normalize credit score
	normalizedCreditScore := (creditScore - MinCreditScore) / (MaxCreditScore - MinCreditScore)

	// Normalize DTI (ensure DTI is not too small)
	if dti < MinDTI {
		dti = MinDTI
	}

	score := (W0 * normalizedCreditScore) + (W1 * (1 / dti))
	return sigmoid(score)
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

/////////////////////// homomorphic

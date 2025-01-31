package credit_evaluation

import (
	"math"
	"testing"
)

// Test satisfyPreselection function
func TestSatisfyPreselection(t *testing.T) {
	tests := []struct {
		name   string
		age    int
		salary int
		want   bool
	}{
		{"Eligible", 25, 150000000, true},
		{"Too Young", 17, 150000000, false},
		{"Salary Too Low", 25, 90000000, false},
		{"Too Young and Salary Too Low", 17, 90000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := satisfyPreselection(tt.age, tt.salary)
			if got != tt.want {
				t.Errorf("satisfyPreselection(%d, %d) = %v, want %v", tt.age, tt.salary, got, tt.want)
			}
		})
	}
}

// Test calcScore function
func TestCalcScore(t *testing.T) {
	tests := []struct {
		name        string
		creditScore float64
		dti         float64
		want        float64
	}{
		{"High Credit Score, Low DTI", 800, 0.2, Sigmoid(W0*(800-MinCreditScore)/(MaxCreditScore-MinCreditScore) + W1*(1/0.2))},
		{"Low Credit Score, High DTI", 500, 0.6, Sigmoid(W0*(500-MinCreditScore)/(MaxCreditScore-MinCreditScore) + W1*(1/0.6))},
		{"Minimum Credit Score, Maximum DTI", 300, 1.0, Sigmoid(W0*(300-MinCreditScore)/(MaxCreditScore-MinCreditScore) + W1*(1/1.0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcScore(tt.creditScore, tt.dti)
			t.Logf("calcScore(%f, %f) = %f", tt.creditScore, tt.dti, got)
			if math.Abs(got-tt.want) > 0.001 { // Allow small floating-point errors
				t.Errorf("calcScore(%f, %f) = %f, want %f", tt.creditScore, tt.dti, got, tt.want)
			}
		})
	}
}

// Test sigmoid function
func TestSigmoid(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"Positive Input", 1.0, 1 / (1 + math.Exp(-1.0))},
		{"Negative Input", -1.0, 1 / (1 + math.Exp(1.0))},
		{"Zero Input", 0.0, 1 / (1 + math.Exp(0.0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sigmoid(tt.x)
			if got != tt.want {
				t.Errorf("sigmoid(%f) = %f, want %f", tt.x, got, tt.want)
			}
		})
	}
}

// Test CreditEvaluation function
func TestCreditEvaluation(t *testing.T) {
	tests := []struct {
		name        string
		age         int
		salary      int
		creditScore float64
		dti         float64
		want        float64
	}{
		{"Eligible and High Score", 25, 150000000, 800, 0.2, calcScore(800, 0.2)},
		{"Eligible and Low Score", 25, 150000000, 500, 0.6, calcScore(500, 0.6)},
		{"Not Eligible (Too Young)", 17, 150000000, 800, 0.2, -1},
		{"Not Eligible (Salary Too Low)", 25, 90000000, 800, 0.2, -1},
		{"Not Eligible (Too Young and Salary Too Low)", 17, 90000000, 800, 0.2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreditEvaluation(tt.age, tt.salary, tt.creditScore, tt.dti)
			if got != tt.want {
				t.Errorf("CreditEvaluation(%d, %d, %f, %f) = %f, want %f", tt.age, tt.salary, tt.creditScore, tt.dti, got, tt.want)
			}
		})
	}
}

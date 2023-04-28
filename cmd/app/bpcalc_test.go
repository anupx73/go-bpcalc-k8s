package main

import "testing"

func TestCalcCategoryLow(t *testing.T) {

	expected := "Low"
	actual := calcCategory(65, 45)

	if actual != expected {
		t.Errorf("Test fail! Expected: '%s', Actual: '%s'", expected, actual)
	}
}

func TestCalcCategoryIdeal(t *testing.T) {

	expected := "Ideal"
	actual := calcCategory(110, 75)

	if actual != expected {
		t.Errorf("Test fail! Expected: '%s', Actual: '%s'", expected, actual)
	}
}

func TestCalcCategoryPreHigh(t *testing.T) {

	expected := "Pre High"
	actual := calcCategory(130, 89)

	if actual != expected {
		t.Errorf("Test fail! Expected: '%s', Actual: '%s'", expected, actual)
	}
}

// func TestCalcCategoryHigh(t *testing.T) {

// 	expected := "High"
// 	actual := calcCategory(180, 99)

// 	if actual != expected {
// 		t.Errorf("Test fail! Expected: '%s', Actual: '%s'", expected, actual)
// 	}
// }

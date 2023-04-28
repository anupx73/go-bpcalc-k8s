package main

// bp category calc
func calcCategory(systolic int, diastolic int) string {
	resCategory := ""

	if systolic < 90 && diastolic < 60 {
		resCategory = "Low"
	} else if systolic < 120 && diastolic < 80 {
		resCategory = "Ideal"
	} else if systolic < 140 && diastolic < 90 {
		resCategory = "Pre High"
	} else if systolic <= 190 && diastolic <= 100 {
		resCategory = "High"
	}

	return resCategory
}

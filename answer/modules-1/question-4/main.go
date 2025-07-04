package main

import "fmt"

type Month uint8

// TODO: Use iota to define the months from January to December.
const (
	_ Month = iota // Ignore 0
	January
	Februari
	March
	April
	Mey
	June
	July
	August
	September
	October
	November
	December
	// ... continue for all 12 months
)

// TODO: Implement this function to convert a Month to its string representation.
func GetMonthName(m Month) string {
	// Hint: A switch statement or a slice of strings would work well here.
	switch m {
	case January:
		return "January"
	case Februari:
		return "Februari"
	case March:
		return "March"
	case April:
		return "April"
	case Mey:
		return "Mey"
	case June:
		return "June"
	case July:
		return "July"
	case August:
		return "August"
	case September:
		return "September"
	case October:
		return "October"
	case November:
		return "November"
	case December:
		return "December"
	default:
		return "Invalid Month"
	}
}

func main() {
	fmt.Println(GetMonthName(January)) // Expected: January
}
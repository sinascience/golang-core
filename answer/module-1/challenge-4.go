package main

import "fmt"

type Month uint8

// TODO: Use iota to define the months from January to December.
const (
	_ Month = iota // Ignore 0
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// TODO: Implement this function to convert a Month to its string representation.
func GetMonthName(m Month) string {
	// Hint: A switch statement or a slice of strings would work well here.
	MonthNames := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if int(m) < 1 || int(m) >= len(MonthNames) {
		return "Unknown Month"
	}
	return MonthNames[m-1]
}

func main() {
	fmt.Println(GetMonthName(January)) // Expected: January
}

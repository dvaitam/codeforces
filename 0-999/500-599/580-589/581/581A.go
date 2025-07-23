package main

import "fmt"

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	var diffDays, sameDays int
	if a < b {
		diffDays = a
		sameDays = (b - a) / 2
	} else {
		diffDays = b
		sameDays = (a - b) / 2
	}
	fmt.Println(diffDays, sameDays)
}

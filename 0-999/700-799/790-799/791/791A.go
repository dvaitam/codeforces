package main

import "fmt"

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	years := 0
	for a <= b {
		a *= 3
		b *= 2
		years++
	}
	fmt.Println(years)
}

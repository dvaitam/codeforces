package main

import (
	"fmt"
)

func main() {
	var n uint64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// For any n >= 2, 5^n ends with 25.
	// The input constraints guarantee n >= 2.
	fmt.Print("25")
}

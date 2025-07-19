package main

import (
	"fmt"
	"os"
)

func main() {
	var n int
	// Read until EOF, handling multiple cases if present
	for {
		if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
			break
		}
		base := "ROYGBIV"
		extra := "GBIV"
		res := base
		// Append extra colors for positions beyond the first 7
		for i := 0; i < n-7; i++ {
			res += string(extra[i%len(extra)])
		}
		fmt.Fprintln(os.Stdout, res)
	}
}

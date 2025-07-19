package main

import (
	"fmt"
	"os"
)

func main() {
	var n int
	var s string
	// Read input
	if _, err := fmt.Fscan(os.Stdin, &n, &s); err != nil {
		return
	}
	// count R and L
	var cR, cL int
	for _, ch := range s {
		if ch == 'R' {
			cR++
		} else if ch == 'L' {
			cL++
		}
	}
	// only L footprints
	if cR == 0 {
		// start: last L
		for i := n - 1; i >= 0; i-- {
			if s[i] == 'L' {
				// print start block
				fmt.Printf("%d ", i+1)
				break
			}
		}
		// end: first L
		for i := 0; i < n; i++ {
			if s[i] == 'L' {
				fmt.Printf("%d", i)
				return
			}
		}
	}
	// only R footprints
	if cL == 0 {
		// start: first R
		for i := 0; i < n; i++ {
			if s[i] == 'R' {
				fmt.Printf("%d ", i+1)
				break
			}
		}
		// end: last R
		for i := n - 1; i >= 0; i-- {
			if s[i] == 'R' {
				fmt.Printf("%d", i+2)
				return
			}
		}
	}
	// both R and L footprints
	// start: first R
	for i := 0; i < n; i++ {
		if s[i] == 'R' {
			fmt.Printf("%d ", i+1)
			break
		}
	}
	// end: first occurrence of RL
	for i := 0; i < n-1; i++ {
		if s[i] == 'R' && s[i+1] == 'L' {
			fmt.Printf("%d", i+1)
			return
		}
	}
}

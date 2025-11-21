package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var s string
		fmt.Fscan(reader, &s)
		solve(s)
	}
}

func solve(s string) {
	n := len(s)
	z := -1
	// Find the first zero
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			z = i
			break
		}
	}

	// If no zeros found (string is all 1s)
	if z == -1 {
		fmt.Printf("1 %d 1 1\n", n)
		return
	}

	// Target length for the second substring
	k := n - z
	bestI := 0
	
	// Iterate through all possible substrings of length k
	// valid indices are 0 to n-k
	for i := 1; i <= n-k; i++ {
		better := false
		// Compare candidate at i with current best candidate at bestI
		// We only need to compare the bits starting from index z, 
		// as the prefix [0...z-1] is always all 1s in the XOR result.
		for j := 0; j < k; j++ {
			// Compare XOR values at the corresponding bit
			valNow := s[z+j] ^ s[i+j]
			valBest := s[z+j] ^ s[bestI+j]
			
			if valNow > valBest {
				better = true
				break
			}
			if valNow < valBest {
				better = false
				break
			}
		}
		
		if better {
			bestI = i
		}
	}

	// Output: l1 r1 l2 r2 (1-based indices)
	fmt.Printf("1 %d %d %d\n", n, bestI+1, bestI+k)
}

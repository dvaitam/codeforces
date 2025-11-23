package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		if isSorted(s) {
			fmt.Fprintln(writer, 0)
			continue
		}

		// 1. Find the lexicographically largest subsequence
		// Indices should be collected by iterating right to left
		// and picking elements >= max_suffix.
		// This gives them in reverse order of indices.
		var subIndices []int
		suffixMax := byte(0)
		for j := n - 1; j >= 0; j-- {
			if s[j] >= suffixMax {
				subIndices = append(subIndices, j)
				suffixMax = s[j]
			}
		}
		// Now reverse subIndices to get them in increasing order
		for j, k := 0, len(subIndices)-1; j < k; j, k = j+1, k-1 {
			subIndices[j], subIndices[k] = subIndices[k], subIndices[j]
		}
		
		// 2. Simulate the cyclic shift.
		// The char at subIndices[i] moves to subIndices[i+1] (cyclically).
		// Right shift means:
		// new char at subIndices[0] = old char at subIndices[len-1]
		// new char at subIndices[1] = old char at subIndices[0]
		// ...
		// new char at subIndices[i] = old char at subIndices[i-1]
		
		sBytes := []byte(s)
		m := len(subIndices)
		
		// Apply the shift
		// Copy original values first
		subVals := make([]byte, m)
		for j := 0; j < m; j++ {
			subVals[j] = s[subIndices[j]]
		}
		
		for j := 0; j < m; j++ {
			// The value at subIndices[j] becomes the value from the *previous* index in the cycle
			prevIdx := (j - 1 + m) % m
			sBytes[subIndices[j]] = subVals[prevIdx]
		}
		
		// 3. Check if sorted
		if !isSorted(string(sBytes)) {
			fmt.Fprintln(writer, -1)
		} else {
			// 4. Calculate cost
			// The subsequence is non-increasing.
			// Sorting it takes (len - count_of_max_element) operations.
			maxVal := subVals[0] // First element is the largest since it's non-increasing
			countMax := 0
			for _, v := range subVals {
				if v == maxVal {
					countMax++
				}
			}
			fmt.Fprintln(writer, m - countMax)
		}
	}
}

func isSorted(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] > s[i+1] {
			return false
		}
	}
	return true
}
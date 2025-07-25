package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func computeMinimal(s string) string {
	// Special case: zero
	if s == "0" {
		return "0"
	}
	// Sort digits ascending
	digits := []rune(s)
	sort.Slice(digits, func(i, j int) bool { return digits[i] < digits[j] })
	// Count leading zeros
	zeroCount := 0
	for zeroCount < len(digits) && digits[zeroCount] == '0' {
		zeroCount++
	}
	// If no leading zeros, join sorted
	if zeroCount == 0 {
		return string(digits)
	}
	// Place first non-zero digit first, then zeros, then rest
	first := digits[zeroCount]
	var b strings.Builder
	b.WriteRune(first)
	for i := 0; i < zeroCount; i++ {
		b.WriteRune('0')
	}
	for i := zeroCount + 1; i < len(digits); i++ {
		b.WriteRune(digits[i])
	}
	return b.String()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// Read original number
	if !scanner.Scan() {
		return
	}
	n := scanner.Text()
	// Read Bob's answer
	if !scanner.Scan() {
		return
	}
	m := scanner.Text()
	// Compute minimal permutation
	correct := computeMinimal(n)
	if m == correct {
		fmt.Println("OK")
	} else {
		fmt.Println("WRONG_ANSWER")
	}
}

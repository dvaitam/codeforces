package main

import (
	"bufio"
	"fmt"
	"os"
)

func parseDigits(s string) []int {
	n := len(s)
	res := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		res[n-1-i] = int(s[i] - '0')
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var A string
	if _, err := fmt.Fscan(in, &A); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	B := make([]string, n)
	maxLen := len(A)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
		if len(B[i]) > maxLen {
			maxLen = len(B[i])
		}
	}
	costs := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Fscan(in, &costs[i])
	}

	// digits of A reversed
	aDigits := parseDigits(A)

	// digits of B reversed and padded
	bDigits := make([][]int, n)
	for i := 0; i < n; i++ {
		bDigits[i] = parseDigits(B[i])
		for len(bDigits[i]) < maxLen+1 {
			bDigits[i] = append(bDigits[i], 0)
		}
	}

	// pad aDigits with zeros
	for len(aDigits) < maxLen {
		aDigits = append(aDigits, 0)
	}

	carry := make([]int, n)
	totalCost := 0

	for pos := 0; pos < len(aDigits); pos++ {
		allowed := []int{}
		if A[len(A)-1-pos] == '?' {
			if pos == len(aDigits)-1 { // most significant digit
				for d := 1; d <= 9; d++ {
					allowed = append(allowed, d)
				}
			} else {
				for d := 0; d <= 9; d++ {
					allowed = append(allowed, d)
				}
			}
		} else {
			allowed = append(allowed, aDigits[pos])
		}

		bestCost := -1
		bestDigit := 0
		bestCarry := make([]int, n)

		for _, dig := range allowed {
			tmpCarry := make([]int, n)
			curCost := 0
			for i := 0; i < n; i++ {
				val := bDigits[i][pos] + carry[i] + dig
				curCost += costs[val%10]
				tmpCarry[i] = val / 10
			}
			if curCost > bestCost {
				bestCost = curCost
				bestDigit = dig
				copy(bestCarry, tmpCarry)
			}
		}
		totalCost += bestCost
		copy(carry, bestCarry)
		aDigits[pos] = bestDigit
	}

	// process remaining positions if needed
	for pos := len(aDigits); pos < maxLen+1; pos++ {
		for i := 0; i < n; i++ {
			val := bDigits[i][pos] + carry[i]
			totalCost += costs[val%10]
			carry[i] = val / 10
		}
	}

	fmt.Fprintln(out, totalCost)
}

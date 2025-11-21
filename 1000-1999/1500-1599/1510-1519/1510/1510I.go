package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func logSum(pred string, target byte, mistakes []int, lnBeta float64) float64 {
	maxLog := math.Inf(-1)
	sum := 0.0
	for i := 0; i < len(pred); i++ {
		if pred[i] != target {
			continue
		}
		lw := float64(mistakes[i]) * lnBeta
		if lw > maxLog {
			if maxLog == math.Inf(-1) {
				sum = 1
			} else {
				sum = sum*math.Exp(maxLog-lw) + 1
			}
			maxLog = lw
		} else {
			sum += math.Exp(lw - maxLog)
		}
	}
	if maxLog == math.Inf(-1) {
		return math.Inf(-1)
	}
	return maxLog + math.Log(sum)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	mistakes := make([]int, n)
	const beta = 0.85
	lnBeta := math.Log(beta)

	for round := 0; round < m; round++ {
		var others string
		if _, err := fmt.Fscan(in, &others); err != nil {
			return
		}
		for len(others) != n {
			// Safety in case of stray whitespace tokens.
			if _, err := fmt.Fscan(in, &others); err != nil {
				return
			}
		}

		log0 := logSum(others, '0', mistakes, lnBeta)
		log1 := logSum(others, '1', mistakes, lnBeta)

		guess := byte('0')
		if log1 > log0 {
			guess = '1'
		}

		fmt.Fprintf(out, "%c\n", guess)
		out.Flush()

		var actual string
		if _, err := fmt.Fscan(in, &actual); err != nil {
			return
		}
		for len(actual) == 0 {
			if _, err := fmt.Fscan(in, &actual); err != nil {
				return
			}
		}
		result := actual[0]

		for i := 0; i < n; i++ {
			if others[i] != result {
				mistakes[i]++
			}
		}
	}
}

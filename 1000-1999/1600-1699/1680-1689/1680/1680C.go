package main

import (
	"bufio"
	"fmt"
	"os"
)

func minCost(s string) int {
	n := len(s)
	zeros := make([]int, n+1)
	ones := make([]int, n+1)
	totalOnes := 0
	for i := 0; i < n; i++ {
		zeros[i+1] = zeros[i]
		ones[i+1] = ones[i]
		if s[i] == '0' {
			zeros[i+1]++
		} else {
			ones[i+1]++
			totalOnes++
		}
	}

	check := func(k int) bool {
		if totalOnes-k <= 0 {
			// No restriction on ones inside the substring
			i := 0
			for j := 0; j <= n; j++ {
				for i < j && zeros[j]-zeros[i] > k {
					i++
				}
				if zeros[j]-zeros[i] <= k {
					return true
				}
			}
			return false
		}
		requiredOnes := totalOnes - k
		i := 0
		for j := 0; j <= n; j++ {
			for i < j && zeros[j]-zeros[i] > k {
				i++
			}
			if ones[j]-ones[i] >= requiredOnes {
				return true
			}
		}
		return false
	}

	l, r := 0, n
	for l < r {
		m := (l + r) / 2
		if check(m) {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, minCost(s))
	}
}

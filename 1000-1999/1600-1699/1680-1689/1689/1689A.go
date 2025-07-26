package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program solves the problem described in problemA.txt for contest 1689.
// It builds the lexicographically smallest string by taking letters from two
// given strings without taking more than k letters from the same string in a row.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		var aStr, bStr string
		fmt.Fscan(reader, &aStr)
		fmt.Fscan(reader, &bStr)

		a := []byte(aStr)
		b := []byte(bStr)
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

		i, j := 0, 0
		cntA, cntB := 0, 0
		var result []byte

		for i < n && j < m {
			if (a[i] < b[j] && cntA < k) || cntB == k {
				result = append(result, a[i])
				i++
				cntA++
				cntB = 0
			} else {
				result = append(result, b[j])
				j++
				cntB++
				cntA = 0
			}
		}
		fmt.Fprintln(writer, string(result))
	}
}

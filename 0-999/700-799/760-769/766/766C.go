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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	a := make([]int, 26)
	for i := 0; i < 26; i++ {
		fmt.Fscan(reader, &a[i])
	}

	const mod int = 1000000007

	dpCount := make([]int, n+1) // number of ways
	dpMin := make([]int, n+1)   // min number of segments
	for i := 1; i <= n; i++ {
		dpMin[i] = 1<<31 - 1
	}
	dpCount[0] = 1
	dpMin[0] = 0
	maxLen := 0

	for i := 1; i <= n; i++ {
		limit := 1000000000
		for j := i; j >= 1; j-- {
			ch := s[j-1] - 'a'
			if a[ch] < limit {
				limit = a[ch]
			}
			length := i - j + 1
			if length > limit {
				break
			}
			if dpCount[j-1] > 0 {
				if length > maxLen {
					maxLen = length
				}
			}
			dpCount[i] = (dpCount[i] + dpCount[j-1]) % mod
			if dpMin[j-1]+1 < dpMin[i] {
				dpMin[i] = dpMin[j-1] + 1
			}
		}
	}

	fmt.Fprintln(writer, dpCount[n])
	fmt.Fprintln(writer, maxLen)
	fmt.Fprintln(writer, dpMin[n])
}

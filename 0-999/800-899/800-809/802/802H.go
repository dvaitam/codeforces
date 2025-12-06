package main

import (
	"fmt"
	"strings"
)

func f(k int) int {
	return k * (k - 1) * (k - 2) * (k - 3) / 24
}

func main() {
	var n int
	fmt.Scan(&n)

	cnt := make([]int, 80)
	k := 79 // Start from a value where f(k) > 10^6 is guaranteed or safe upper bound.
            // f(80) approx 1.5e6 > 1e6. So 79 is safe start.

	for n > 0 {
		for f(k) > n {
			k--
		}
		n -= f(k)
		cnt[k]++
	}

	var sb strings.Builder
	for i := 1; i < 80; i++ {
		sb.WriteByte('a')
		for j := 0; j < cnt[i]; j++ {
			sb.WriteByte('b')
		}
	}
	sb.WriteString(" aaaab")
	fmt.Println(sb.String())
}
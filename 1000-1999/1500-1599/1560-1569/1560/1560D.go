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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	powers := precompute()

	for ; t > 0; t-- {
		var n string
		fmt.Fscan(reader, &n)
		minOps := len(n) + 1e9 // big number
		for _, p := range powers {
			l := lcsPrefix(n, p)
			ops := len(n) - l + len(p) - l
			if ops < minOps {
				minOps = ops
			}
		}
		fmt.Fprintln(writer, minOps)
	}
}

// precompute returns decimal strings of powers of two up to 2^60.
func precompute() []string {
	res := make([]string, 0, 63)
	for i := 0; i <= 60; i++ {
		res = append(res, fmt.Sprintf("%d", 1<<uint(i)))
	}
	return res
}

// lcsPrefix returns the length of the longest prefix of b that is a subsequence of a.
func lcsPrefix(a, b string) int {
	j := 0
	for i := 0; i < len(a) && j < len(b); i++ {
		if a[i] == b[j] {
			j++
		}
	}
	return j
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		b := []byte(s)
		sorted := make([]byte, n)
		copy(sorted, b)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

		diff := 0
		for i := 0; i < n; i++ {
			if b[i] != sorted[i] {
				diff++
			}
		}
		fmt.Fprintln(out, diff)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		bytes := []byte(s)
		sort.Slice(bytes, func(i, j int) bool { return bytes[i] < bytes[j] })

		if bytes[0] != bytes[k-1] {
			fmt.Fprintf(out, "%c\n", bytes[k-1])
			continue
		}

		// start result with first character
		res := []byte{bytes[0]}
		if k < n {
			if bytes[k] != bytes[n-1] {
				res = append(res, bytes[k:]...)
			} else {
				// all remaining characters equal
				count := (n - k + k - 1) / k
				res = append(res, []byte(strings.Repeat(string(bytes[k]), count))...)
			}
		}
		fmt.Fprintln(out, string(res))
	}
}

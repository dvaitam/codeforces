package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It counts the number of integers k that satisfy a set of
// constraints on being greater than, less than, or not equal
// to given values.
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
		lo := 0
		hi := int(1e9)
		banned := make(map[int]struct{})
		for i := 0; i < n; i++ {
			var a, x int
			fmt.Fscan(in, &a, &x)
			switch a {
			case 1:
				if x > lo {
					lo = x
				}
			case 2:
				if x < hi {
					hi = x
				}
			case 3:
				banned[x] = struct{}{}
			}
		}
		if lo > hi {
			fmt.Fprintln(out, 0)
			continue
		}
		count := hi - lo + 1
		for x := range banned {
			if x >= lo && x <= hi {
				count--
			}
		}
		if count < 0 {
			count = 0
		}
		fmt.Fprintln(out, count)
	}
}

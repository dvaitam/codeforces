package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the offline version of problemG3: given a permutation p of
// [1..n], output the index y such that p[y] = 1.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	ans := -1
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x == 1 {
			ans = i
		}
	}
	if ans != -1 {
		fmt.Fprintln(out, ans)
	}
}

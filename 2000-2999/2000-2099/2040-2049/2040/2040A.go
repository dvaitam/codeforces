package main

import (
	"bufio"
	"fmt"
	"os"
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
		a := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			rem := ((a[i] % k) + k) % k
			freq[rem]++
		}
		ansIdx := -1
		for i := 0; i < n; i++ {
			rem := ((a[i] % k) + k) % k
			if freq[rem] == 1 {
				ansIdx = i + 1
				break
			}
		}
		if ansIdx == -1 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, ansIdx)
		}
	}
}


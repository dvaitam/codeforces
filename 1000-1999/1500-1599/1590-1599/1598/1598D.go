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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		cntA := make(map[int]int)
		cntB := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
			cntA[a[i]]++
			cntB[b[i]]++
		}
		total := int64(n) * int64(n-1) * int64(n-2) / 6
		var bad int64
		for i := 0; i < n; i++ {
			bad += int64(cntA[a[i]]-1) * int64(cntB[b[i]]-1)
		}
		fmt.Fprintln(out, total-bad)
	}
}

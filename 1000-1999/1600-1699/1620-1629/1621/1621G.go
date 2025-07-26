package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

func weight(a []int, idx []int) int {
	if len(idx) == 0 {
		return 0
	}
	last := idx[len(idx)-1]
	maxAfter := 0
	for i := last + 1; i < len(a); i++ {
		if a[i] > maxAfter {
			maxAfter = a[i]
		}
	}
	w := 0
	for _, p := range idx {
		if a[p] < maxAfter {
			w++
		}
	}
	return w % MOD
}

func dfs(a []int, start int, subseq []int, ans *int) {
	if len(subseq) > 0 {
		*ans = (*ans + weight(a, subseq)) % MOD
	}
	for i := start; i < len(a); i++ {
		if len(subseq) == 0 || a[i] > a[subseq[len(subseq)-1]] {
			subseq = append(subseq, i)
			dfs(a, i+1, subseq, ans)
			subseq = subseq[:len(subseq)-1]
		}
	}
}

func solve(a []int) int {
	ans := 0
	dfs(a, 0, []int{}, &ans)
	return ans % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		fmt.Fprintln(writer, solve(a))
	}
}

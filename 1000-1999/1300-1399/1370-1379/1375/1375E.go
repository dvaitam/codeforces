package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ u, v int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	ops := make([]pair, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if a[i] > a[j] {
				ops = append(ops, pair{i, j})
			}
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		ai, aj := a[ops[i].u], a[ops[j].u]
		if ai != aj {
			return ai < aj
		}
		return ops[i].v > ops[j].v
	})
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, len(ops))
	for _, p := range ops {
		fmt.Fprintln(out, p.u+1, p.v+1)
	}
}

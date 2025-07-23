package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	i, j int
	v    int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	m := 2 * n
	pairs := make([]Pair, 0, m*(m-1)/2)
	// strengths are given for i from 2..m
	for i := 2; i <= m; i++ {
		for j := 1; j < i; j++ {
			var val int
			fmt.Fscan(reader, &val)
			pairs = append(pairs, Pair{i, j, val})
		}
	}

	sort.Slice(pairs, func(a, b int) bool {
		return pairs[a].v > pairs[b].v
	})

	res := make([]int, m+1)
	for _, p := range pairs {
		if res[p.i] == 0 && res[p.j] == 0 {
			res[p.i] = p.j
			res[p.j] = p.i
		}
	}

	for i := 1; i <= m; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
}

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

	var n int
	fmt.Fscan(in, &n)
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}

	res := make([]int, n)
	for k := 1; k <= n; k++ {
		res[k-1] = countSegments(colors, k)
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}

func countSegments(a []int, k int) int {
	if k == 0 {
		return len(a)
	}
	segments := 1
	freq := make(map[int]int)
	distinct := 0
	for _, c := range a {
		if freq[c] == 0 {
			if distinct == k {
				segments++
				freq = map[int]int{c: 1}
				distinct = 1
				continue
			}
			distinct++
		}
		freq[c]++
	}
	return segments
}

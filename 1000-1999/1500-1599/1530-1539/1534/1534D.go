package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(out *bufio.Writer, in *bufio.Reader, v, n int) []int {
	fmt.Fprintf(out, "? %d\n", v)
	out.Flush()
	d := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &d[i])
	}
	return d
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	dist1 := query(out, in, 1, n)
	even := []int{}
	odd := []int{}

	edges := make(map[[2]int]struct{})
	for i := 2; i <= n; i++ {
		if dist1[i] == 1 {
			a, b := 1, i
			if a > b {
				a, b = b, a
			}
			edges[[2]int{a, b}] = struct{}{}
		}
		if dist1[i]%2 == 0 {
			even = append(even, i)
		} else {
			odd = append(odd, i)
		}
	}

	group := even
	if len(odd) < len(even) {
		group = odd
	}

	for _, v := range group {
		if v == 1 {
			continue
		}
		d := query(out, in, v, n)
		for i := 1; i <= n; i++ {
			if d[i] == 1 {
				a, b := v, i
				if a > b {
					a, b = b, a
				}
				edges[[2]int{a, b}] = struct{}{}
			}
		}
	}

	fmt.Fprintln(out, "!")
	for e := range edges {
		fmt.Fprintf(out, "%d %d\n", e[0], e[1])
	}
}

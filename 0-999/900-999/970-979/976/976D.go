package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v int
}

// construct builds edges for degree set d starting from vertex index st.
// The algorithm follows the constructive proof from the editorial.
func construct(st int, d []int) []edge {
	if len(d) == 0 {
		return nil
	}

	res := make([]edge, 0)
	last := d[len(d)-1]
	for i := 0; i < d[0]; i++ {
		for j := st + i + 1; j <= st+last; j++ {
			res = append(res, edge{st + i, j})
		}
	}

	// prepare reduced degree set for recursive call
	nxt := st + d[0]
	if len(d) > 1 {
		nd := make([]int, len(d)-1)
		for i := 1; i < len(d); i++ {
			nd[i-1] = d[i] - d[0]
		}
		nd = nd[:len(nd)-1]
		res = append(res, construct(nxt, nd)...)
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	d := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &d[i])
	}

	edges := construct(0, d)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, len(edges))
	for _, e := range edges {
		fmt.Fprintf(out, "%d %d\n", e.u+1, e.v+1)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var tt int
	if _, err := fmt.Fscan(reader, &tt); err != nil {
		return
	}
	for tt > 0 {
		tt--
		var n, m int
		fmt.Fscan(reader, &n, &m)
		// build reverse graph: from b to a
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			a--
			b--
			g[b] = append(g[b], a)
		}
		// BFS from node 0
		d := make([]int, n)
		for i := range d {
			d[i] = -1
		}
		queue := make([]int, 0, n)
		d[0] = 1
		queue = append(queue, 0)
		for qi := 0; qi < len(queue); qi++ {
			u := queue[qi]
			for _, v := range g[u] {
				if d[v] == -1 {
					d[v] = d[u] + 1
					queue = append(queue, v)
				}
			}
		}
		// check unreachable
		infinite := false
		for i := 0; i < n; i++ {
			if d[i] == -1 {
				infinite = true
				break
			}
		}
		if infinite {
			fmt.Fprintln(writer, "INFINITE")
			continue
		}
		fmt.Fprintln(writer, "FINITE")
		// bucket nodes by distance
		at := make([][]int, n+1)
		for i := 0; i < n; i++ {
			di := d[i]
			if di >= 0 && di <= n {
				at[di] = append(at[di], i)
			}
		}
		// build sequence
		// total capacity sum of distances
		total := 0
		for i := 1; i <= n; i++ {
			total += len(at[i])
		}
		// but nodes appear multiple times: sum(di)
		capSeq := 0
		for i := 0; i < n; i++ {
			capSeq += d[i]
		}
		seq := make([]int, 0, capSeq)
		for from := 1; from <= n; from++ {
			for val := n; val >= from; val-- {
				for _, x := range at[val] {
					seq = append(seq, x)
				}
			}
		}
		// output
		fmt.Fprintln(writer, len(seq))
		for i, x := range seq {
			// print 1-indexed
			if i+1 < len(seq) {
				fmt.Fprintf(writer, "%d ", x+1)
			} else {
				fmt.Fprintf(writer, "%d\n", x+1)
			}
		}
	}
}

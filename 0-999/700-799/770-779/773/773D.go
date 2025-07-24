package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	w := make([][]int64, n)
	for i := range w {
		w[i] = make([]int64, n)
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			var x int64
			fmt.Fscan(reader, &x)
			w[i][j] = x
			w[j][i] = x
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	const inf int64 = 1 << 60
	for r := 0; r < n; r++ {
		d := make([]int64, n)
		vis := make([]bool, n)
		for i := 0; i < n; i++ {
			if i == r {
				d[i] = inf
			} else {
				d[i] = w[i][r]
			}
		}
		vis[r] = true
		var sum int64
		for cnt := 0; cnt < n-1; cnt++ {
			u := -1
			best := inf
			for i := 0; i < n; i++ {
				if !vis[i] && d[i] < best {
					best = d[i]
					u = i
				}
			}
			vis[u] = true
			sum += d[u]
			for v := 0; v < n; v++ {
				if !vis[v] {
					cand := d[u]
					if w[u][v] < cand {
						cand = w[u][v]
					}
					if cand < d[v] {
						d[v] = cand
					}
				}
			}
		}
		if r > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, sum)
	}
	writer.WriteByte('\n')
}

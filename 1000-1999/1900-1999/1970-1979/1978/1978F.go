package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		// construct matrix b with cyclic shifts
		b := make([][]int, n)
		for i := 0; i < n; i++ {
			b[i] = make([]int, n)
			for j := 0; j < n; j++ {
				idx := (j - i + n) % n
				b[i][j] = a[idx]
			}
		}

		total := n * n
		visited := make([]bool, total)
		comps := 0
		for idx := 0; idx < total; idx++ {
			if visited[idx] {
				continue
			}
			comps++
			queue := []int{idx}
			visited[idx] = true
			for len(queue) > 0 {
				v := queue[0]
				queue = queue[1:]
				x1 := v / n
				y1 := v % n
				for w := 0; w < total; w++ {
					if visited[w] {
						continue
					}
					x2 := w / n
					y2 := w % n
					if abs(x1-x2)+abs(y1-y2) <= k && gcd(b[x1][y1], b[x2][y2]) > 1 {
						visited[w] = true
						queue = append(queue, w)
					}
				}
			}
		}
		fmt.Fprintln(out, comps)
	}
}

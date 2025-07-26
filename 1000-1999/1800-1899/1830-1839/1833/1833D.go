package main

import (
	"bufio"
	"fmt"
	"os"
)

func lexGreater(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if n == 1 {
			fmt.Fprintln(out, p[0])
			continue
		}
		// positions of maximum excluding first and last elements
		posMax := 1
		for i := 1; i < n; i++ {
			if p[i] > p[posMax] {
				posMax = i
			}
		}
		posPrevMax := 0
		for i := 0; i < n-1; i++ {
			if p[i] > p[posPrevMax] {
				posPrevMax = i
			}
		}
		rsetMap := map[int]bool{}
		candidates := []int{posMax - 1, posMax, posPrevMax - 1, posPrevMax, n - 1}
		for _, r := range candidates {
			if r >= 0 && r < n {
				rsetMap[r] = true
			}
		}
		best := make([]int, 0)
		for r := range rsetMap {
			suffix := make([]int, n-r-1)
			copy(suffix, p[r+1:])
			for l := 0; l <= r; l++ {
				cand := make([]int, 0, n)
				cand = append(cand, suffix...)
				for i := r; i >= l; i-- {
					cand = append(cand, p[i])
				}
				cand = append(cand, p[:l]...)
				if len(best) == 0 || lexGreater(cand, best) {
					best = cand
				}
			}
		}
		for i, val := range best {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, val)
		}
		out.WriteByte('\n')
	}
}

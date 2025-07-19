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

	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	type seg struct{ l, r, id int }
	segs := make([]seg, m)
	for i := 0; i < m; i++ {
		// Convert to 0-based
		fmt.Fscan(reader, &segs[i].l, &segs[i].r)
		segs[i].l--
		segs[i].r--
		segs[i].id = i + 1
	}
	// Handle special cases
	if m == 0 && n >= 2 {
		mx, mi := a[0], a[0]
		for _, v := range a {
			if v > mx {
				mx = v
			}
			if v < mi {
				mi = v
			}
		}
		fmt.Fprintln(writer, mx-mi)
		fmt.Fprintln(writer, 0)
		return
	}
	if n == 1 {
		fmt.Fprintln(writer, 0)
		fmt.Fprintln(writer, 0)
		return
	}
	bestDiff := -1 << 60
	var bestSegs []int
	// Try removing by each position i
	for i := 0; i < n; i++ {
		// copy a to b
		b := make([]int, n)
		copy(b, a)
		var cur []int
		for _, s := range segs {
			if s.l <= i && i <= s.r {
				cur = append(cur, s.id)
				for k := s.l; k <= s.r; k++ {
					b[k]--
				}
			}
		}
		mx, mi := b[0], b[0]
		for _, v := range b {
			if v > mx {
				mx = v
			}
			if v < mi {
				mi = v
			}
		}
		diff := mx - mi
		if diff > bestDiff {
			bestDiff = diff
			bestSegs = make([]int, len(cur))
			copy(bestSegs, cur)
		}
	}
	fmt.Fprintln(writer, bestDiff)
	fmt.Fprintln(writer, len(bestSegs))
	for i, id := range bestSegs {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, id)
	}
	if len(bestSegs) > 0 {
		writer.WriteByte('\n')
	}
}

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type segment struct {
	x1, y1, x2, y2 float64
	length         float64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	segs := make([]segment, n)
	total := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		x1, y1 := xs[i], ys[i]
		x2, y2 := xs[j], ys[j]
		l := math.Abs(x1-x2) + math.Abs(y1-y2)
		segs[i] = segment{x1, y1, x2, y2, l}
		total += l
	}

	delta := total / float64(m)

	prefix := make([]float64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + segs[i].length
	}

	pos := func(t float64) (float64, float64) {
		t = math.Mod(t, total)
		if t < 0 {
			t += total
		}
		idx := sort.Search(len(prefix)-1, func(i int) bool { return prefix[i+1] > t })
		seg := segs[idx]
		rem := t - prefix[idx]
		if seg.x1 == seg.x2 {
			if seg.y2 > seg.y1 {
				return seg.x1, seg.y1 + rem
			}
			return seg.x1, seg.y1 - rem
		}
		if seg.x2 > seg.x1 {
			return seg.x1 + rem, seg.y1
		}
		return seg.x1 - rem, seg.y1
	}

	evaluate := func(start float64) float64 {
		maxD := 0.0
		for i := 0; i < m; i++ {
			x1, y1 := pos(start + float64(i)*delta)
			x2, y2 := pos(start + float64(i+1)*delta)
			d := math.Hypot(x1-x2, y1-y2)
			if d > maxD {
				maxD = d
			}
		}
		return maxD
	}

	steps := n * 4
	if steps < 200 {
		steps = 200
	}
	best := math.MaxFloat64
	for i := 0; i < steps; i++ {
		s := total * float64(i) / float64(steps)
		d := evaluate(s)
		if d < best {
			best = d
		}
	}

	fmt.Fprintf(out, "%.10f\n", best)
}

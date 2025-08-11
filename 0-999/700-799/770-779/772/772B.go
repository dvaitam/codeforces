package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type pt struct{ x, y float64 }

func pointLineDistance(p, a, b pt) float64 {
	num := math.Abs((p.x-a.x)*(b.y-a.y) - (p.y-a.y)*(b.x-a.x))
	den := math.Hypot(b.x-a.x, b.y-a.y)
	if den == 0 {
		return math.Hypot(p.x-a.x, p.y-a.y)
	}
	return num / den
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	p := make([]pt, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i].x, &p[i].y)
	}

	ans := math.MaxFloat64
	for i := 0; i < n; i++ {
		a := p[(i-1+n)%n]
		b := p[i]
		c := p[(i+1)%n]

		altA := pointLineDistance(a, b, c)
		altB := pointLineDistance(b, a, c)
		altC := pointLineDistance(c, a, b)

		d := math.Min(altA, math.Min(altB, altC)) / 2
		if d < ans {
			ans = d
		}
	}
	fmt.Fprintf(out, "%.10f\n", ans)
}

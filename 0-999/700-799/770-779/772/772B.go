package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type pt struct{ x, y float64 }
	p := make([]pt, n)
	for i := 0; i < n; i++ {
		var xi, yi float64
		fmt.Fscan(in, &xi, &yi)
		p[i] = pt{xi, yi}
	}

	ans := math.MaxFloat64
	for i := 0; i < n; i++ {
		a := p[(i-1+n)%n]
		b := p[i]
		c := p[(i+1)%n]
		cross := math.Abs((b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x))
		base := math.Hypot(c.x-a.x, c.y-a.y)
		h := cross / base
		d := h / 2
		if d < ans {
			ans = d
		}
	}
	fmt.Fprintf(out, "%.10f\n", ans)
}

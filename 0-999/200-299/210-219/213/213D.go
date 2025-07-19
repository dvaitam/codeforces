package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Vec holds a 2D point or vector with components y and x
type Vec struct {
	y, x float64
}

// add returns the vector sum of v and r
func (v Vec) add(r Vec) Vec {
	return Vec{v.y + r.y, v.x + r.x}
}

// sub returns the vector difference of v and r
func (v Vec) sub(r Vec) Vec {
	return Vec{v.y - r.y, v.x - r.x}
}

// rotate rotates vector l by angle r (radians)
func rotate(l Vec, r float64) Vec {
	return Vec{
		l.y*math.Cos(r) + l.x*math.Sin(r),
		l.x*math.Cos(r) - l.y*math.Sin(r),
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	// precompute constants
	PI := math.Atan2(0, -1)
	a := 72.0 / 180.0 * PI
	base := Vec{0, 10}

	// build points
	pts := make([]Vec, 0, 1+4*n)
	pts = append(pts, Vec{0, 0})
	for i := 0; i < n; i++ {
		c := len(pts) - 1
		pts = append(pts, pts[c].sub(rotate(base, -2*a)))
		pts = append(pts, pts[c].add(rotate(base, -a)))
		pts = append(pts, pts[len(pts)-1].add(base))
		pts = append(pts, pts[len(pts)-1].add(rotate(base, a)))
	}

	// output points
	fmt.Fprintln(out, len(pts))
	for _, v := range pts {
		fmt.Fprintf(out, "%.9f %.9f\n", v.x, v.y)
	}

	// output faces/sequences
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "%d %d %d %d %d\n", i*4+1, i*4+3, i*4+4, i*4+5, i*4+2)
	}
	// special sequence
	fmt.Fprint(out, 1)
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, " %d %d %d %d", i*4+4, i*4+2, i*4+3, i*4+5)
	}
	for i := n - 1; i >= 0; i-- {
		fmt.Fprintf(out, " %d", i*4+1)
	}
	fmt.Fprintln(out)
}

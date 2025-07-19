package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

const eps = 1e-10

// vect represents a 2D vector or point
type vect struct{ x, y float64 }

func (v vect) add(o vect) vect    { return vect{v.x + o.x, v.y + o.y} }
func (v vect) sub(o vect) vect    { return vect{v.x - o.x, v.y - o.y} }
func (v vect) mul(s float64) vect { return vect{v.x * s, v.y * s} }
func (v vect) len() float64       { return math.Hypot(v.x, v.y) }
func (v vect) lensq() float64     { return v.x*v.x + v.y*v.y }

// circle with center v and radius r
type circle struct {
	v vect
	r float64
}

// isin reports whether circle a is inside circle b
func isin(a, b circle) bool {
	return a.v.sub(b.v).len() <= b.r-a.r+eps
}

// f2 returns minimal circle enclosing two circles a and b
func f2(a, b circle) circle {
	v := b.v.sub(a.v)
	d := v.len()
	// compute weighted mid
	va := v.mul(a.r / d)
	vb := v.mul((d - b.r) / d)
	ctr := a.v.add(va.add(vb).mul(0.5))
	rad := va.sub(vb).len() * 0.5
	return circle{ctr, rad}
}

// f3 returns minimal circle enclosing three circles a, b, c
func f3(a, b, c circle) circle {
	x1, y1, r1 := a.v.x, a.v.y, a.r
	x2, y2, r2 := b.v.x, b.v.y, b.r
	x3, y3, r3 := c.v.x, c.v.y, c.r
	a2 := x1 - x2
	a3 := x1 - x3
	b2 := y1 - y2
	b3 := y1 - y3
	c2 := r2 - r1
	c3 := r3 - r1
	d1 := x1*x1 + y1*y1 - r1*r1
	d2 := d1 - x2*x2 - y2*y2 + r2*r2
	d3 := d1 - x3*x3 - y3*y3 + r3*r3
	ab := a3*b2 - a2*b3
	xa := (b2*d3-b3*d2)/(2*ab) - x1
	xb := (b3*c2 - b2*c3) / ab
	ya := (a3*d2-a2*d3)/(2*ab) - y1
	yb := (a2*c3 - a3*c2) / ab
	A := xb*xb + yb*yb - 1
	B := 2 * (r1 + xa*xb + ya*yb)
	C := xa*xa + ya*ya - r1*r1
	var r float64
	if math.Abs(A) > eps {
		disc := B*B - 4*A*C
		r = -(B - math.Sqrt(disc)) / (2 * A)
	} else {
		r = -C / B
	}
	cx := x1 + xa + xb*r
	cy := y1 + ya + yb*r
	return circle{vect{cx, cy}, r}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	v := make([]circle, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &v[i].v.x, &v[i].v.y, &v[i].r)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { v[i], v[j] = v[j], v[i] })
	ans := circle{vect{0, 0}, 1e18}
   for i := 0; i < n; i++ {
       if !isin(ans, v[i]) {
           ans = v[i]
           for j := 0; j < i; j++ {
               if !isin(ans, v[j]) {
                   ans = f2(v[i], v[j])
                   for k := 0; k < j; k++ {
                       if !isin(ans, v[k]) {
                           ans = f3(v[i], v[j], v[k])
                       }
                   }
               }
           }
       }
   }
	// output with high precision
	fmt.Printf("%.20f %.20f %.20f", ans.v.x, ans.v.y, ans.r)
}

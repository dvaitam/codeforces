package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// PT represents a point or vector in 2D.
type PT struct {
	x, y float64
}

// Add returns the vector sum of p and q.
func (p PT) Add(q PT) PT { return PT{p.x + q.x, p.y + q.y} }

// Sub returns the vector difference p - q.
func (p PT) Sub(q PT) PT { return PT{p.x - q.x, p.y - q.y} }

// Mul returns vector p scaled by t.
func (p PT) Mul(t float64) PT { return PT{p.x * t, p.y * t} }

// Div returns vector p divided by scalar t.
func (p PT) Div(t float64) PT { return PT{p.x / t, p.y / t} }

// Rot returns p rotated by angle t (radians) counter-clockwise.
func (p PT) Rot(t float64) PT {
	c, s := math.Cos(t), math.Sin(t)
	return PT{p.x*c - p.y*s, p.x*s + p.y*c}
}

// Dot returns the dot product of p and q.
func Dot(p, q PT) float64 { return p.x*q.x + p.y*q.y }

// Cross returns the cross product (scalar) of p and q.
func Cross(p, q PT) float64 { return p.x*q.y - p.y*q.x }

// Angle returns the signed angle from p to q in radians.
func Angle(p, q PT) float64 { return math.Atan2(Cross(p, q), Dot(p, q)) }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	pts := make([]PT, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}
	// Compute centroid
	var cen PT
	var mass float64
	for i := 2; i < n; i++ {
		p0 := pts[0]
		p1 := pts[i-1]
		p2 := pts[i]
		temp := PT{(p0.x + p1.x + p2.x) / 3.0, (p0.y + p1.y + p2.y) / 3.0}
		area2 := math.Abs(Cross(p1.Sub(p0), p2.Sub(p0)))
		cen = cen.Mul(mass).Add(temp.Mul(area2)).Div(mass + area2)
		mass += area2
	}
	// Translate points to center
	for i := 0; i < n; i++ {
		pts[i] = pts[i].Sub(cen)
	}
	// Process queries
	a, b := 0, 1
	ang := 0.0
	const twoPi = 2 * math.Pi
	up := PT{0, 1}
	for i := 0; i < q; i++ {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var c1, c2 int
			fmt.Fscan(in, &c1)
			c1--
			if b == c1 {
				a, b = b, a
			}
			// Rotate about pts[b]
			r := pts[b].Rot(ang)
			cen = cen.Add(r)
			// Compute rotation to align r with up vector
			tang := Angle(r, up)
			ang += tang
			// normalize angle
			ang = math.Mod(ang, twoPi)
			if ang < 0 {
				ang += twoPi
			}
			// update center
			cen = cen.Sub(pts[b].Rot(ang))
			fmt.Fscan(in, &c2)
			a = c2 - 1
		} else {
			var c int
			fmt.Fscan(in, &c)
			c--
			r := pts[c].Rot(ang)
			p := r.Add(cen)
			fmt.Fprintf(out, "%.8f %.8f\n", p.x, p.y)
		}
	}
}

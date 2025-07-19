package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// P represents a point or vector in 3D space
type P struct{ x, y, z float64 }

func (p P) add(q P) P       { return P{p.x + q.x, p.y + q.y, p.z + q.z} }
func (p P) sub(q P) P       { return P{p.x - q.x, p.y - q.y, p.z - q.z} }
func (p P) mul(k float64) P { return P{p.x * k, p.y * k, p.z * k} }
func (p P) dot(q P) float64 { return p.x*q.x + p.y*q.y + p.z*q.z }
func (p P) cross(q P) P {
	return P{
		p.y*q.z - p.z*q.y,
		p.z*q.x - p.x*q.z,
		p.x*q.y - p.y*q.x,
	}
}
func (p P) mag2() float64 { return p.dot(p) }
func (p P) mag() float64  { return math.Sqrt(p.mag2()) }

const (
	eps = 1e-9
	inf = 1e20
)

// f checks intersection time between moving sphere (centered at a, radius r1, velocity v)
// and fixed sphere (centered at o, radius r2). Returns (true, t) if t>0 solution exists.
func f(a, o P, r1, r2 float64, v P) (bool, float64) {
	A := v.mag2()
	B := 2 * a.sub(o).dot(v)
	C := a.sub(o).mag2() - (r1+r2)*(r1+r2)
	D := B*B - 4*A*C
	if D < 0 {
		return false, 0
	}
	if D > 0 {
		D = math.Sqrt(D)
	}
	t := (-B - D) / (2 * A)
	return t > 0, t
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var A, V P
	var R float64
	var n int
	fmt.Fscan(in, &A.x, &A.y, &A.z)
	fmt.Fscan(in, &V.x, &V.y, &V.z)
	fmt.Fscan(in, &R)
	fmt.Fscan(in, &n)
	ans := inf
	upd := func(t float64) {
		if t < ans {
			ans = t
		}
	}
	for i := 0; i < n; i++ {
		var cen P
		var r float64
		var m int
		fmt.Fscan(in, &cen.x, &cen.y, &cen.z)
		fmt.Fscan(in, &r, &m)
		if ok, t := f(A, cen, R, r, V); ok {
			upd(t)
		}
		for j := 0; j < m; j++ {
			var p P
			fmt.Fscan(in, &p.x, &p.y, &p.z)
			if ok, t := f(A, cen, R, 0, V); ok {
				upd(t)
			}
			if ok, t := f(A, cen.add(p), R, 0, V); ok {
				upd(t)
			}
			c := cen
			d := cen.add(p)
			cd := d.sub(c)
			ca := c.sub(A)
			crossVc := V.cross(cd)
			A1 := crossVc.mag2()
			if math.Abs(A1) < eps {
				continue
			}
			B := crossVc.mag() * ca.cross(cd).mag()
			C := ca.cross(cd).mag2()
			D := B*B - 4*A1*C
			if D >= 0 {
				if D > 0 {
					D = math.Sqrt(D)
				}
				t1 := (-B + D) / (2 * A1)
				if t1 > 0 {
					Q := A.add(V.mul(t1))
					if Q.sub(c).dot(cd) >= 0 && Q.sub(d).dot(c.sub(d)) >= 0 {
						upd(t1)
					}
				}
			}
		}
	}
	if ans == inf {
		fmt.Println(-1)
	} else {
		fmt.Printf("%.20f\n", ans)
	}
}

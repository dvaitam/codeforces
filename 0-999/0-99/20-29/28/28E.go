package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const eps = 1e-13

var n int
var xCo, yCo []int64

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func vect(x1, y1, x2, y2 int64) int64 {
	return x1*y2 - x2*y1
}

// point-in-polygon exact for integer coords
func inLL(xx, yy int64) bool {
	var sum, rl int64
	for i := 0; i < n; i++ {
		sum += abs64((xCo[i]-xx)*(yCo[i+1]-yy) - (xCo[i+1]-xx)*(yCo[i]-yy))
		rl += xCo[i]*yCo[i+1] - xCo[i+1]*yCo[i]
	}
	return sum == abs64(rl)
}

// point-in-polygon with floats for approximate check
func in2(xx, yy float64) bool {
	var sum, rl float64
	for i := 0; i < n; i++ {
		xi, yi := float64(xCo[i]), float64(yCo[i])
		xi1, yi1 := float64(xCo[i+1]), float64(yCo[i+1])
		sum += math.Abs((xi-xx)*(yi1-yy) - (xi1-xx)*(yi-yy))
		rl += xi*yi1 - xi1*yi
	}
	return math.Abs(sum-math.Abs(rl)) <= eps
}

// intersection of line segment and ray for integer coords
func cross(x1, y1, x2, y2, x3, y3, x4, y4 int64, r1, r2 int) (bool, float64) {
	a := x2 - x1
	b := x3 - x4
	c := x3 - x1
	d := y2 - y1
	e := y3 - y4
	f := y3 - y1
	Q := a*e - b*d
	if Q != 0 {
		T := c*e - b*f
		S := a*f - c*d
		if Q < 0 {
			T, S, Q = -T, -S, -Q
		}
		ok := true
		if r1&1 != 0 && T < 0 {
			ok = false
		}
		if r1&2 != 0 && T > Q {
			ok = false
		}
		if r2&1 != 0 && S < 0 {
			ok = false
		}
		if r2&2 != 0 && S > Q {
			ok = false
		}
		if ok {
			return true, float64(T) / float64(Q)
		}
	}
	return false, 0
}

// intersection with floats for approximate check
func cross2(x1, y1, x2, y2, x3, y3, x4, y4 float64, r1, r2 int) (bool, float64) {
	a := x2 - x1
	b := x3 - x4
	c := x3 - x1
	d := y2 - y1
	e := y3 - y4
	f := y3 - y1
	Q := a*e - b*d
	if Q != 0 {
		T := c*e - b*f
		S := a*f - c*d
		if Q < 0 {
			T, S, Q = -T, -S, -Q
		}
		ok := true
		if r1&1 != 0 && T < -eps*Q {
			ok = false
		}
		if r1&2 != 0 && T > Q*(1+eps) {
			ok = false
		}
		if r2&1 != 0 && S < -eps*Q {
			ok = false
		}
		if r2&2 != 0 && S > Q*(1+eps) {
			ok = false
		}
		if ok {
			return true, T / Q
		}
	}
	return false, 0
}

// check for integer line
func check(x1, y1, x2, y2 int64, r1, r2 int) (bool, float64) {
	if inLL(x1, y1) {
		return true, 0
	}
	ok := false
	tMin := 1e20
	for i := 0; i < n; i++ {
		if ok2, t := cross(0, 0, x1, y1, xCo[i], yCo[i], xCo[i+1], yCo[i+1], r1, r2); ok2 {
			if t < tMin {
				tMin = t
			}
			ok = true
		}
	}
	return ok, tMin
}

// check with floats
func check2(x1, y1, x2, y2 float64, r1, r2 int) (bool, float64) {
	if in2(x1, y1) {
		return true, 0
	}
	ok := false
	tMin := 1e20
	for i := 0; i < n; i++ {
		xi, yi := float64(xCo[i]), float64(yCo[i])
		xi1, yi1 := float64(xCo[i+1]), float64(yCo[i+1])
		if ok2, t := cross2(x1, y1, x2, y2, xi, yi, xi1, yi1, r1, r2); ok2 {
			if t < tMin {
				tMin = t
			}
			ok = true
		}
	}
	return ok, tMin
}

func main() {
	rdr := bufio.NewReader(os.Stdin)
	wrt := bufio.NewWriter(os.Stdout)
	defer wrt.Flush()
	// read input
	var X, Y int64
	fmt.Fscan(rdr, &n)
	xCo = make([]int64, n+1)
	yCo = make([]int64, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &xCo[i], &yCo[i])
	}
	fmt.Fscan(rdr, &X, &Y)
	for i := 0; i < n; i++ {
		xCo[i] -= X
		yCo[i] -= Y
	}
	// close polygon
	xCo[n] = xCo[0]
	yCo[n] = yCo[0]
	var x1, y1, z1, f, x2, y2, z2 int64
	fmt.Fscan(rdr, &x1, &y1, &z1)
	fmt.Fscan(rdr, &f)
	f = -f
	fmt.Fscan(rdr, &x2, &y2, &z2)
	z2 = -z2
	x3 := x1*z2 + x2*z1
	y3 := y1*z2 + y2*z1
	// search times
	t1, t2 := 1e20, 1e20
	// first direct drop
	if ok, t := check(0, 0, x1, y1, 1, 3); ok {
		_, tt := check2(float64(x1)*t+float64(x2)*(float64(z1)*t/float64(z2)),
			float64(y1)*t+float64(y2)*(float64(z1)*t/float64(z2)),
			float64(x1)*t, float64(y1)*t, 3, 3)
		tt *= (float64(z1) * t / float64(f))
		if t1 > t+eps || (math.Abs(t-t1) < eps && t2 > tt+eps) {
			t1 = t
			t2 = tt
		}
	}
	// via x3,y3
	if ok, t := check(0, 0, x3, y3, 1, 3); ok {
		t = t * float64(z2)
		_, tt := check2(float64(x1)*t+float64(x2)*(float64(z1)*t/float64(z2)),
			float64(y1)*t+float64(y2)*(float64(z1)*t/float64(z2)),
			float64(x1)*t, float64(y1)*t, 3, 3)
		tt *= (float64(z1) * t / float64(f))
		if t1 > t+eps || (math.Abs(t-t1) < eps && t2 > tt+eps) {
			t1 = t
			t2 = tt
		}
	}
	// edge case
	if vect(x1, y1, x2, y2) != 0 {
		T := 1e20
		for i := 0; i < n; i++ {
			if ok2, r := cross(0, 0, x1, y1,
				xCo[i], yCo[i], xCo[i]-x2*z1, yCo[i]-y2*z1, 1, 1); ok2 {
				if r < T {
					T = r
				}
			}
		}
		if T < 5e19 {
			t := T
			if ok2, tt := check2(float64(x1)*T+float64(x2)*(float64(z1)*T/float64(z2)),
				float64(y1)*T+float64(y2)*(float64(z1)*T/float64(z2)),
				float64(x1)*T, float64(y1)*T, 3, 3); ok2 {
				tt *= (float64(z1) * T / float64(f))
				if t1 > t+eps || (math.Abs(t-t1) < eps && t2 > tt+eps) {
					t1 = t
					t2 = tt
				}
			}
		}
	}
	// output
	if t1 > 5e19 {
		fmt.Fprintln(wrt, "-1 -1")
	} else {
		fmt.Fprintf(wrt, "%.15f %.15f\n", t1, t2)
	}
}

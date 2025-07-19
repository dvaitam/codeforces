package main

import (
	"bufio"
	"fmt"
	"os"
)

// check returns true if during movement from (X1,Y1) to (X1+a1, Y1+b1)
// and (X2,Y2) to (X2+a2, Y2+b2) the distance between points <= sqrt(d1sq).
func check(X1, Y1, X2, Y2, a1, b1, a2, b2, d1sq int64) bool {
	dx := a1 - a2
	dy := b1 - b2
	// Quadratic coefficients for distance squared between two moving points
	// f(t) = A*t^2 + B*t + C
	A := dx*dx + dy*dy
	B := 2 * (dx*(X1-X2) + dy*(Y1-Y2))
	C := (X1-X2)*(X1-X2) + (Y1-Y2)*(Y1-Y2)
	// at t=0
	if C <= d1sq {
		return true
	}
	// at t=1
	if A+B+C <= d1sq {
		return true
	}
	// if extremum inside (0,1)
	if A != 0 && B <= 0 && -B <= 2*A {
		// check minimum value <= d1sq: 4*A*C - B*B <= 4*d1sq*A
		if 4*A*C-B*B <= 4*d1sq*A {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var d1, d2 int64
	if _, err := fmt.Fscan(reader, &n, &d1, &d2); err != nil {
		return
	}
	d1sq := d1 * d1
	d2sq := d2 * d2

	// read initial positions
	var ax, ay, bx, by int64
	fmt.Fscan(reader, &ax, &ay, &bx, &by)
	prevAx, prevAy, prevBx, prevBy := ax, ay, bx, by
	// initial distance
	dx0 := ax - bx
	dy0 := ay - by
	ans := 0
	flag := false
	if dx0*dx0+dy0*dy0 <= d1sq {
		flag = true
		ans++
	}
	// process movements
	for i := 1; i < n; i++ {
		fmt.Fscan(reader, &ax, &ay, &bx, &by)
		da1 := ax - prevAx
		db1 := ay - prevAy
		da2 := bx - prevBx
		db2 := by - prevBy
		f := check(prevAx, prevAy, prevBx, prevBy, da1, db1, da2, db2, d1sq)
		if !flag && f {
			ans++
		}
		if f {
			flag = true
		}
		// update flag if out of d2 range
		ddx := ax - bx
		ddy := ay - by
		if ddx*ddx+ddy*ddy > d2sq {
			flag = false
		}
		prevAx, prevAy, prevBx, prevBy = ax, ay, bx, by
	}
	fmt.Fprintln(writer, ans)
}

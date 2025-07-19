package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	// prepare geometry
	pi := math.Pi
	fn := float64(n)
	dirX := math.Cos(pi / fn)
	dirY := math.Sin(pi / fn)
	// points A, B, C
	// A = (0,0)
	// B = (1,0)
	cX, cY := -dirX, dirY
	// B is (1,0)
	// for each segment
	type pair struct {
		val float64
		idx int
	}
	pp := make([]pair, n)
	for i := 0; i < n; i++ {
		x := a[i]
		y := a[(i+1)%n]
		// E = C + (A-C)*x/n = C*(1 - x/n)
		fx := float64(x) / fn
		ex := cX * (1 - fx)
		ey := cY * (1 - fx)
		// F = B + (D-B)*y/n = (1,0) + (dirX,dirY)*y/n
		fy := float64(y) / fn
		fx2 := dirX * fy
		fy2 := dirY * fy
		// F point coordinates
		fxp := 1 + fx2
		fyp := fy2
		// vector FE = F - E
		ux := fxp - ex
		uy := fyp - ey
		// area(E,F,A) = |det(FE, A-E)|, A-E = (-ex, -ey)
		det1 := ux*(-ey) - uy*(-ex)
		a1 := math.Abs(det1)
		// area(E,F,B) = |det(FE, B-E)|, B-E = (1-ex, -ey)
		det2 := ux*(-ey) - uy*(1-ex)
		a2 := math.Abs(det2)
		pp[i] = pair{a1 - a2, i}
	}
	sort.Slice(pp, func(i, j int) bool {
		return pp[i].val < pp[j].val
	})
	ans := make([]int, n)
	// assign ranks: smallest gets n-1, largest gets 0
	for i, p := range pp {
		ans[p.idx] = n - 1 - i
	}
	// output
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, ans[i])
	}
	out.WriteByte('\n')
}

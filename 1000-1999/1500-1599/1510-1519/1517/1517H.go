package main

import (
	"bufio"
	"fmt"
	"os"
)

// interval represents [l, r] range inclusive
// All values fit in int64

type interval struct {
	l int64
	r int64
}

func intersect(it *interval, l, r int64) bool {
	if l > it.l {
		it.l = l
	}
	if r < it.r {
		it.r = r
	}
	return it.l <= it.r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		xL := make([]int64, n+1)
		xR := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &xL[i], &xR[i])
		}
		yL := make([]int64, n+1)
		yR := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(reader, &yL[i], &yR[i])
		}
		zL := make([]int64, n+1)
		zR := make([]int64, n+1)
		for i := 3; i <= n; i++ {
			fmt.Fscan(reader, &zL[i], &zR[i])
		}

		b := make([]interval, n+1)
		for i := 1; i <= n; i++ {
			b[i] = interval{l: xL[i], r: xR[i]}
		}
		d := make([]interval, n+1) // d[i] = b[i]-b[i-1] for i>=2
		for i := 2; i <= n; i++ {
			d[i] = interval{l: yL[i], r: yR[i]}
		}

		feasible := true
		for iter := 0; iter < n*3 && feasible; iter++ {
			changed := false
			// equality constraints between b and d
			for i := 2; i <= n && feasible; i++ {
				// update d from b
				if !intersect(&d[i], b[i].l-b[i-1].r, b[i].r-b[i-1].l) {
					feasible = false
					break
				}
			}
			for i := 2; i <= n && feasible; i++ {
				// update b from d
				if !intersect(&b[i], b[i-1].l+d[i].l, b[i-1].r+d[i].r) {
					feasible = false
					break
				}
				if !intersect(&b[i-1], b[i].l-d[i].r, b[i].r-d[i].l) {
					feasible = false
					break
				}
			}
			// apply y bounds explicitly
			for i := 2; i <= n && feasible; i++ {
				if !intersect(&d[i], yL[i], yR[i]) {
					feasible = false
					break
				}
			}
			// second difference constraints
			for i := 3; i <= n && feasible; i++ {
				if !intersect(&d[i], d[i-1].l+zL[i], d[i-1].r+zR[i]) {
					feasible = false
					break
				}
				if !intersect(&d[i-1], d[i].l-zR[i], d[i].r-zL[i]) {
					feasible = false
					break
				}
			}
			if !feasible {
				break
			}
			// ensure b within original bounds
			for i := 1; i <= n; i++ {
				nl, nr := b[i].l, b[i].r
				if nl < xL[i] {
					nl = xL[i]
				}
				if nr > xR[i] {
					nr = xR[i]
				}
				if nl != b[i].l || nr != b[i].r {
					changed = true
					b[i].l, b[i].r = nl, nr
					if nl > nr {
						feasible = false
						break
					}
				}
			}
			for i := 2; i <= n && feasible; i++ {
				nl, nr := d[i].l, d[i].r
				if nl < yL[i] {
					nl = yL[i]
				}
				if nr > yR[i] {
					nr = yR[i]
				}
				if nl != d[i].l || nr != d[i].r {
					changed = true
					d[i].l, d[i].r = nl, nr
					if nl > nr {
						feasible = false
						break
					}
				}
			}
			if !changed {
				break
			}
		}
		if !feasible {
			fmt.Fprintln(writer, "NO")
			continue
		}
		// final verification
		ok := true
		for i := 1; i <= n; i++ {
			if b[i].l > b[i].r {
				ok = false
				break
			}
		}
		if ok {
			for i := 2; i <= n; i++ {
				diffL := b[i].l - b[i-1].r
				diffR := b[i].r - b[i-1].l
				if diffL > yR[i] || diffR < yL[i] {
					ok = false
					break
				}
			}
		}
		if ok {
			for i := 3; i <= n; i++ {
				low := (b[i].l - b[i-1].r) - (b[i-1].r - b[i-2].l)
				high := (b[i].r - b[i-1].l) - (b[i-1].l - b[i-2].r)
				if low > zR[i] || high < zL[i] {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

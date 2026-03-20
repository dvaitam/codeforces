package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	ID        int
	X, Y      int
	IsBerland bool
	R         int
}

func crossProduct(A, B, C Point) int64 {
	return int64(B.X-A.X)*int64(C.Y-A.Y) - int64(B.Y-A.Y)*int64(C.X-A.X)
}

func getConvexHull(pts []Point) []Point {
	sorted := make([]Point, len(pts))
	copy(sorted, pts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].X != sorted[j].X {
			return sorted[i].X < sorted[j].X
		}
		return sorted[i].Y < sorted[j].Y
	})

	var hull []Point
	for _, p := range sorted {
		for len(hull) >= 2 && crossProduct(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}

	t := len(hull) + 1
	for i := len(sorted) - 2; i >= 0; i-- {
		p := sorted[i]
		for len(hull) >= t && crossProduct(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}

	if len(hull) > 1 {
		hull = hull[:len(hull)-1]
	}
	return hull
}

func getRestSorted(pts []Point, C Point) []Point {
	var rest []Point
	for _, p := range pts {
		if p.ID != C.ID || p.IsBerland != C.IsBerland {
			rest = append(rest, p)
		}
	}
	sort.Slice(rest, func(i, j int) bool {
		return crossProduct(C, rest[i], rest[j]) > 0
	})
	return rest
}

func solve(pts []Point, out *bufio.Writer) {
	if len(pts) <= 1 {
		return
	}
	if len(pts) == 2 {
		if pts[0].IsBerland {
			fmt.Fprintf(out, "%d %d\n", pts[0].ID, pts[1].ID)
		} else {
			fmt.Fprintf(out, "%d %d\n", pts[1].ID, pts[0].ID)
		}
		return
	}

	hull := getConvexHull(pts)

	for _, p := range hull {
		if p.IsBerland && p.R >= 2 {
			rest := getRestSorted(pts, p)
			S := 0
			k := -1
			for i := 0; i < len(rest); i++ {
				if rest[i].IsBerland {
					S += rest[i].R - 1
				} else {
					S -= 1
				}
				if S == -1 {
					k = i
					break
				}
			}
			fmt.Fprintf(out, "%d %d\n", p.ID, rest[k].ID)
			L := rest[:k+1]
			R := rest[k+1:]
			if len(L) > 0 {
				solve(L, out)
			}
			pNew := p
			pNew.R--
			if pNew.R > 0 {
				rNew := append([]Point{}, R...)
				rNew = append(rNew, pNew)
				solve(rNew, out)
			}
			return
		}
	}

	for _, p := range hull {
		if p.IsBerland {
			rest := getRestSorted(pts, p)
			S := 0
			k := -1
			for i := 0; i < len(rest); i++ {
				if rest[i].IsBerland {
					S += rest[i].R - 1
				} else {
					S -= 1
				}
				if S == -1 {
					k = i
					break
				}
			}
			if k == len(rest)-1 {
				fmt.Fprintf(out, "%d %d\n", p.ID, rest[k].ID)
				L := rest[:k+1]
				solve(L, out)
				return
			}
		} else {
			rest := getRestSorted(pts, p)
			S := 0
			k := -1
			S_prev := 0
			for i := 0; i < len(rest); i++ {
				S_next := S
				if rest[i].IsBerland {
					S_next += rest[i].R - 1
				} else {
					S_next -= 1
				}
				if S <= -1 && S_next >= 1 {
					k = i
					S_prev = S
					break
				}
				S = S_next
			}
			if k != -1 {
				Pk := rest[k]
				fmt.Fprintf(out, "%d %d\n", Pk.ID, p.ID)
				dL := -1 - S_prev
				dR := S_prev + Pk.R - 2

				L := rest[:k]
				if dL+1 > 0 {
					PkL := Pk
					PkL.R = dL + 1
					LNew := append([]Point{}, L...)
					LNew = append(LNew, PkL)
					solve(LNew, out)
				}

				R := rest[k+1:]
				if dR+1 > 0 {
					PkR := Pk
					PkR.R = dR + 1
					RNew := append([]Point{}, R...)
					RNew = append(RNew, PkR)
					solve(RNew, out)
				}
				return
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for tc := 0; tc < t; tc++ {
		var a, b int
		fmt.Fscan(in, &a, &b)

		r := make([]int, b)
		for i := 0; i < b; i++ {
			fmt.Fscan(in, &r[i])
		}

		pts := make([]Point, 0, a+b)
		for i := 1; i <= a; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			pts = append(pts, Point{ID: i, X: x, Y: y, IsBerland: false, R: 0})
		}
		for i := 1; i <= b; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			pts = append(pts, Point{ID: i, X: x, Y: y, IsBerland: true, R: r[i-1]})
		}

		fmt.Fprintln(out, "YES")
		solve(pts, out)
	}
}

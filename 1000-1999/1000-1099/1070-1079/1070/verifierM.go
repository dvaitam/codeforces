package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// refSolve is the embedded reference solver from cf_latest_1070_M.go
func refSolve(input []byte) string {
	in := bufio.NewReader(bytes.NewReader(input))
	var out bytes.Buffer
	w := bufio.NewWriter(&out)

	type Point struct {
		ID        int
		X, Y      int
		IsBerland bool
		R         int
	}

	crossProduct := func(A, B, C Point) int64 {
		return int64(B.X-A.X)*int64(C.Y-A.Y) - int64(B.Y-A.Y)*int64(C.X-A.X)
	}

	getConvexHull := func(pts []Point) []Point {
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

	getRestSorted := func(pts []Point, C Point) []Point {
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

	var solve func(pts []Point, w *bufio.Writer)
	solve = func(pts []Point, w *bufio.Writer) {
		if len(pts) <= 1 {
			return
		}
		if len(pts) == 2 {
			if pts[0].IsBerland {
				fmt.Fprintf(w, "%d %d\n", pts[0].ID, pts[1].ID)
			} else {
				fmt.Fprintf(w, "%d %d\n", pts[1].ID, pts[0].ID)
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
				fmt.Fprintf(w, "%d %d\n", p.ID, rest[k].ID)
				L := rest[:k+1]
				R := rest[k+1:]
				if len(L) > 0 {
					solve(L, w)
				}
				pNew := p
				pNew.R--
				if pNew.R > 0 {
					rNew := append([]Point{}, R...)
					rNew = append(rNew, pNew)
					solve(rNew, w)
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
					fmt.Fprintf(w, "%d %d\n", p.ID, rest[k].ID)
					L := rest[:k+1]
					solve(L, w)
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
					fmt.Fprintf(w, "%d %d\n", Pk.ID, p.ID)
					dL := -1 - S_prev
					dR := S_prev + Pk.R - 2

					L := rest[:k]
					if dL+1 > 0 {
						PkL := Pk
						PkL.R = dL + 1
						LNew := append([]Point{}, L...)
						LNew = append(LNew, PkL)
						solve(LNew, w)
					}

					R := rest[k+1:]
					if dR+1 > 0 {
						PkR := Pk
						PkR.R = dR + 1
						RNew := append([]Point{}, R...)
						RNew = append(RNew, PkR)
						solve(RNew, w)
					}
					return
				}
			}
		}
	}

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return ""
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

		fmt.Fprintln(w, "YES")
		solve(pts, w)
	}
	w.Flush()
	return strings.TrimSpace(out.String())
}

func runProg(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) []byte {
	t := rng.Intn(2) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for ; t > 0; t-- {
		a := 2
		b := 1
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		sb.WriteString("2\n")
		for i := 0; i < a; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(10), rng.Intn(10)))
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(10), rng.Intn(10)))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		want := refSolve(input)
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

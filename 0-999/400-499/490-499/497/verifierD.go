package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type IPoint struct{ x, y int64 }

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var Px, Py, Qx, Qy int64
	var n, m int
	if _, err := fmt.Fscan(in, &Px, &Py); err != nil {
		return ""
	}
	fmt.Fscan(in, &n)
	A := make([]IPoint, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i].x, &A[i].y)
	}
	fmt.Fscan(in, &Qx, &Qy)
	fmt.Fscan(in, &m)
	B := make([]IPoint, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &B[i].x, &B[i].y)
	}
	Dx := Qx - Px
	Dy := Qy - Py
	Cx := float64(-Dx)
	Cy := float64(-Dy)
	R2i := Dx*Dx + Dy*Dy
	if R2i == 0 {
		return "NO\n"
	}
	R2 := float64(R2i)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dx := A[i].x - B[j].x + Dx
			dy := A[i].y - B[j].y + Dy
			if dx*dx+dy*dy == R2i {
				return "YES\n"
			}
		}
	}
	eps := 1e-9
	for j := 0; j < m; j++ {
		bvx := float64(B[j].x)
		bvy := float64(B[j].y)
		for i := 0; i < n; i++ {
			ni := (i + 1) % n
			u1x := float64(A[i].x) - bvx
			u1y := float64(A[i].y) - bvy
			u2x := float64(A[ni].x) - bvx
			u2y := float64(A[ni].y) - bvy
			dx := u2x - u1x
			dy := u2y - u1y
			a := dx*dx + dy*dy
			ux := u1x - Cx
			uy := u1y - Cy
			if a < eps {
				if math.Abs(ux*ux+uy*uy-R2) < eps {
					return "YES\n"
				}
				continue
			}
			bq := 2 * (dx*ux + dy*uy)
			cq := ux*ux + uy*uy - R2
			disc := bq*bq - 4*a*cq
			if disc < 0 {
				continue
			}
			sd := math.Sqrt(disc)
			t1 := (-bq - sd) / (2 * a)
			t2 := (-bq + sd) / (2 * a)
			if (t1 >= -eps && t1 <= 1+eps) || (t2 >= -eps && t2 <= 1+eps) {
				return "YES\n"
			}
		}
	}
	for i := 0; i < n; i++ {
		avx := float64(A[i].x)
		avy := float64(A[i].y)
		for j := 0; j < m; j++ {
			nj := (j + 1) % m
			u1x := avx - float64(B[j].x)
			u1y := avy - float64(B[j].y)
			u2x := avx - float64(B[nj].x)
			u2y := avy - float64(B[nj].y)
			dx := u2x - u1x
			dy := u2y - u1y
			a := dx*dx + dy*dy
			ux := u1x - Cx
			uy := u1y - Cy
			if a < eps {
				if math.Abs(ux*ux+uy*uy-R2) < eps {
					return "YES\n"
				}
				continue
			}
			bq := 2 * (dx*ux + dy*uy)
			cq := ux*ux + uy*uy - R2
			disc := bq*bq - 4*a*cq
			if disc < 0 {
				continue
			}
			sd := math.Sqrt(disc)
			t1 := (-bq - sd) / (2 * a)
			t2 := (-bq + sd) / (2 * a)
			if (t1 >= -eps && t1 <= 1+eps) || (t2 >= -eps && t2 <= 1+eps) {
				return "YES\n"
			}
		}
	}
	return "NO\n"
}

func randPoly(rng *rand.Rand, n int) []IPoint {
	pts := make([]IPoint, n)
	for i := 0; i < n; i++ {
		pts[i] = IPoint{int64(rng.Intn(21) - 10), int64(rng.Intn(21) - 10)}
	}
	cx, cy := 0.0, 0.0
	for _, p := range pts {
		cx += float64(p.x)
		cy += float64(p.y)
	}
	cx /= float64(n)
	cy /= float64(n)
	sort.Slice(pts, func(i, j int) bool {
		ai := math.Atan2(float64(pts[i].y)-cy, float64(pts[i].x)-cx)
		aj := math.Atan2(float64(pts[j].y)-cy, float64(pts[j].x)-cx)
		return ai < aj
	})
	return pts
}

func genCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	P := IPoint{int64(rng.Intn(11) - 5), int64(rng.Intn(11) - 5)}
	Q := IPoint{int64(rng.Intn(11) - 5), int64(rng.Intn(11) - 5)}
	if P.x == Q.x && P.y == Q.y {
		Q.x++
	}
	A := randPoly(rng, n)
	B := randPoly(rng, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", P.x, P.y)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", A[i].x, A[i].y)
	}
	fmt.Fprintf(&sb, "%d %d\n", Q.x, Q.y)
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", B[i].x, B[i].y)
	}
	inp := sb.String()
	exp := solveD(inp)
	return inp, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

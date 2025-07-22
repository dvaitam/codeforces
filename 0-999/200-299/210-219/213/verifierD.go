package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Vec struct{ y, x float64 }

func (v Vec) add(r Vec) Vec { return Vec{v.y + r.y, v.x + r.x} }
func (v Vec) sub(r Vec) Vec { return Vec{v.y - r.y, v.x - r.x} }
func rotate(l Vec, r float64) Vec {
	return Vec{l.y*math.Cos(r) + l.x*math.Sin(r), l.x*math.Cos(r) - l.y*math.Sin(r)}
}

func genCaseD(rng *rand.Rand) int {
	return rng.Intn(5) + 1
}

func solveD(n int) string {
	PI := math.Atan2(0, -1)
	a := 72.0 / 180.0 * PI
	base := Vec{0, 10}
	pts := make([]Vec, 0, 1+4*n)
	pts = append(pts, Vec{0, 0})
	for i := 0; i < n; i++ {
		c := len(pts) - 1
		pts = append(pts, pts[c].sub(rotate(base, -2*a)))
		pts = append(pts, pts[c].add(rotate(base, -a)))
		pts = append(pts, pts[len(pts)-1].add(base))
		pts = append(pts, pts[len(pts)-1].add(rotate(base, a)))
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, len(pts))
	for _, v := range pts {
		fmt.Fprintf(&sb, "%.9f %.9f\n", v.x, v.y)
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", i*4+1, i*4+3, i*4+4, i*4+5, i*4+2)
	}
	fmt.Fprint(&sb, 1)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, " %d %d %d %d", i*4+4, i*4+2, i*4+3, i*4+5)
	}
	for i := n - 1; i >= 0; i-- {
		fmt.Fprintf(&sb, " %d", i*4+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runD(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	got := strings.TrimSpace(string(out))
	exp := strings.TrimSpace(solveD(n))
	if got != exp {
		return fmt.Errorf("output mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		n := genCaseD(rng)
		if err := runD(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

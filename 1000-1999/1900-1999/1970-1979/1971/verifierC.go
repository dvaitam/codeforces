package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func coord(p int) (float64, float64) {
	if p == 12 {
		p = 0
	}
	angle := 2 * math.Pi * float64(p) / 12
	return math.Cos(angle), math.Sin(angle)
}

func orientation(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

func intersect(a, b, c, d int) bool {
	ax, ay := coord(a)
	bx, by := coord(b)
	cx, cy := coord(c)
	dx, dy := coord(d)
	o1 := orientation(ax, ay, bx, by, cx, cy)
	o2 := orientation(ax, ay, bx, by, dx, dy)
	o3 := orientation(cx, cy, dx, dy, ax, ay)
	o4 := orientation(cx, cy, dx, dy, bx, by)
	return o1*o2 < 0 && o3*o4 < 0
}

type caseC struct{ a, b, c, d int }

func generateCase(rng *rand.Rand) caseC {
	vals := rng.Perm(12)[:4]
	for i := range vals {
		vals[i]++
	}
	return caseC{vals[0], vals[1], vals[2], vals[3]}
}

func runCase(bin string, tc caseC) error {
	input := fmt.Sprintf("1\n%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.ToLower(strings.TrimSpace(out.String()))
	exp := "no"
	if intersect(tc.a, tc.b, tc.c, tc.d) {
		exp = "yes"
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%d %d %d %d\n", i+1, err, tc.a, tc.b, tc.c, tc.d)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

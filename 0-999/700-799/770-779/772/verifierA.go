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

type testCase struct {
	n int
	p float64
	a []float64
	b []float64
}

func feasible(mid float64, a, b []float64, p float64) bool {
	var need float64
	for i := range a {
		use := mid * a[i]
		if use > b[i] {
			need += use - b[i]
		}
	}
	return need <= p*mid
}

func solveCase(tc testCase) string {
	sumA := 0.0
	for _, v := range tc.a {
		sumA += v
	}
	if sumA <= tc.p {
		return "-1\n"
	}
	l, r := 0.0, 1e18
	for i := 0; i < 100; i++ {
		mid := (l + r) / 2
		if feasible(mid, tc.a, tc.b, tc.p) {
			l = mid
		} else {
			r = mid
		}
	}
	return fmt.Sprintf("%.6f\n", l)
}

func genCase(rng *rand.Rand) (string, string) {
	tc := testCase{}
	tc.n = rng.Intn(8) + 1
	tc.p = float64(rng.Intn(10) + 1)
	tc.a = make([]float64, tc.n)
	tc.b = make([]float64, tc.n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, int(tc.p))
	for i := 0; i < tc.n; i++ {
		ai := float64(rng.Intn(10) + 1)
		bi := float64(rng.Intn(10) + 1)
		tc.a[i] = ai
		tc.b[i] = bi
		fmt.Fprintf(&sb, "%d %d\n", int(ai), int(bi))
	}
	exp := solveCase(tc)
	return sb.String(), exp
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	exp = strings.TrimSpace(exp)
	if exp == "-1" {
		if got != "-1" {
			return fmt.Errorf("expected -1 got %s", got)
		}
		return nil
	}
	var g, e float64
	fmt.Sscan(got, &g)
	fmt.Sscan(exp, &e)
	if math.Abs(g-e) > 1e-4*math.Max(1, math.Abs(e)) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

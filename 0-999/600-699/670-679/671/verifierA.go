package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseA struct {
	ax, ay float64
	bx, by float64
	tx, ty float64
	n      int
	pts    [][2]float64
}

func genTestsA() []testCaseA {
	rand.Seed(1)
	tests := make([]testCaseA, 100)
	for i := range tests {
		tc := testCaseA{}
		tc.ax = float64(rand.Intn(6))
		tc.ay = float64(rand.Intn(6))
		tc.bx = float64(rand.Intn(6))
		tc.by = float64(rand.Intn(6))
		tc.tx = float64(rand.Intn(6))
		tc.ty = float64(rand.Intn(6))
		tc.n = rand.Intn(4) + 1 // 1..4
		tc.pts = make([][2]float64, tc.n)
		for j := 0; j < tc.n; j++ {
			tc.pts[j][0] = float64(rand.Intn(6))
			tc.pts[j][1] = float64(rand.Intn(6))
		}
		tests[i] = tc
	}
	return tests
}

func solveA(tc testCaseA) float64 {
	ax, ay := tc.ax, tc.ay
	bx, by := tc.bx, tc.by
	tx, ty := tc.tx, tc.ty
	n := tc.n
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = tc.pts[i][0]
		y[i] = tc.pts[i][1]
	}
	sum := 0.0
	distT := make([]float64, n)
	for i := 0; i < n; i++ {
		d := math.Hypot(x[i]-tx, y[i]-ty)
		distT[i] = d
		sum += 2 * d
	}
	ans := math.Inf(1)
	for i := 0; i < n; i++ {
		costA := math.Hypot(x[i]-ax, y[i]-ay)
		cand := costA + sum - distT[i]
		if cand < ans {
			ans = cand
		}
		costB := math.Hypot(x[i]-bx, y[i]-by)
		cand = costB + sum - distT[i]
		if cand < ans {
			ans = cand
		}
	}
	if n == 1 {
		return ans
	}
	deltaA := make([]float64, n)
	deltaB := make([]float64, n)
	for i := 0; i < n; i++ {
		deltaA[i] = math.Hypot(x[i]-ax, y[i]-ay) - distT[i]
		deltaB[i] = math.Hypot(x[i]-bx, y[i]-by) - distT[i]
	}
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return deltaB[idx[i]] < deltaB[idx[j]] })
	first := idx[0]
	second := idx[1%len(idx)]
	for i := 0; i < n; i++ {
		total := sum + deltaA[i]
		if i == first {
			total += deltaB[second]
		} else {
			total += deltaB[first]
		}
		if total < ans {
			ans = total
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%.0f %.0f %.0f %.0f %.0f %.0f %d\n", tc.ax, tc.ay, tc.bx, tc.by, tc.tx, tc.ty, tc.n)
		for j := 0; j < tc.n; j++ {
			fmt.Fprintf(&input, "%.0f %.0f\n", tc.pts[j][0], tc.pts[j][1])
		}
		expect := solveA(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if math.Abs(val-expect) > 1e-6 {
			fmt.Fprintf(os.Stderr, "test %d: expected %.6f got %.6f\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

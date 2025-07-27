package main

import (
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

type point struct{ x, y int }

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildOracle() (string, error) {
	tmp, err := os.CreateTemp("", "oracle*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), "1381E.go").CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n       int
	q       int
	poly    []point
	queries []int
}

func genPolygon() []point {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rng.Intn(4) + 3
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	poly := make([]point, n)
	for i := 0; i < n; i++ {
		r := rng.Float64()*5 + 5
		x := int(math.Round(r * math.Cos(angles[i]) * 10))
		y := int(math.Round(r * math.Sin(angles[i]) * 10))
		if i > 0 && y == poly[i-1].y {
			y++
		}
		poly[i] = point{x, y}
	}
	return poly
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	basePoly := []point{{0, 0}, {0, 10}, {10, 9}, {10, -1}}
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		poly := basePoly
		n := len(poly)
		q := rng.Intn(3) + 1
		qs := make([]int, 0, q)
		used := map[int]bool{}
		for len(qs) < q {
			f := rng.Intn(9) + 1
			if !used[f] {
				used[f] = true
				qs = append(qs, f)
			}
		}
		tests[i] = testCase{n, q, poly, qs}
	}
	return tests
}

func parseFloats(s string, cnt int) ([]float64, error) {
	r := strings.NewReader(s)
	vals := make([]float64, cnt)
	for i := 0; i < cnt; i++ {
		if _, err := fmt.Fscan(r, &vals[i]); err != nil {
			return nil, err
		}
	}
	return vals, nil
}

func compare(a, b float64) bool {
	diff := math.Abs(a - b)
	den := math.Max(1.0, math.Abs(b))
	return diff/den <= 1e-4
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for _, p := range tc.poly {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
		}
		for _, f := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d\n", f))
		}
		input := sb.String()
		wantStr, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle fail case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		gotStr, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		wantVals, err := parseFloats(wantStr, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output parse error case %d: %v", i+1, err)
			os.Exit(1)
		}
		gotVals, err := parseFloats(gotStr, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\noutput:%s", i+1, err, gotStr)
			os.Exit(1)
		}
		for j := 0; j < tc.q; j++ {
			if !compare(wantVals[j], gotVals[j]) {
				fmt.Fprintf(os.Stderr, "case %d mismatch on query %d\nexpected:%f got:%f\ninput:%s", i+1, j+1, wantVals[j], gotVals[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

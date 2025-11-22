package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type sensorsData struct {
	l [][]int // y, z
	r [][]int // y, z
	f [][]int // x, z
	b [][]int // x, z
	d [][]int // x, y
	u [][]int // x, y
}

type testCase struct {
	n, m, k int
	s       sensorsData
}

type grid struct {
	n, m, k int
	cells   []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for i, tc := range tests {
		input := serializeInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		refGrid, refImpossible, err := parseOutput(refOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\noutput:\n%s", i+1, err, refOut)
			os.Exit(1)
		}
		if !refImpossible {
			if err := validateGrid(tc, refGrid); err != nil {
				fmt.Fprintf(os.Stderr, "reference produced invalid grid on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
		}

		candOut, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candGrid, candImpossible, err := parseOutput(candOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, candOut)
			os.Exit(1)
		}

		if refImpossible != candImpossible {
			fmt.Fprintf(os.Stderr, "test %d: expected impossible=%v, candidate impossible=%v\ninput:\n%s", i+1, refImpossible, candImpossible, input)
			os.Exit(1)
		}
		if candImpossible {
			continue
		}
		if err := validateGrid(tc, candGrid); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate grid invalid: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-1320F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref1320F")
	cmd := exec.Command("go", "build", "-o", outPath, "1320F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func serializeInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n\n", tc.n, tc.m, tc.k))

	writeMatrix := func(mat [][]int) {
		for i := 0; i < len(mat); i++ {
			for j := 0; j < len(mat[i]); j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(mat[i][j]))
			}
			sb.WriteByte('\n')
		}
	}

	writeMatrix(tc.s.l)
	sb.WriteByte('\n')
	writeMatrix(tc.s.r)
	sb.WriteByte('\n')
	writeMatrix(tc.s.f)
	sb.WriteByte('\n')
	writeMatrix(tc.s.b)
	sb.WriteByte('\n')
	writeMatrix(tc.s.d)
	sb.WriteByte('\n')
	writeMatrix(tc.s.u)

	return sb.String()
}

func parseOutput(out string, tc testCase) (grid, bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return grid{}, false, fmt.Errorf("empty output")
	}
	if len(tokens) == 1 && tokens[0] == "-1" {
		return grid{}, true, nil
	}
	expect := tc.n * tc.m * tc.k
	if len(tokens) != expect {
		return grid{}, false, fmt.Errorf("expected %d integers, got %d", expect, len(tokens))
	}
	g := grid{n: tc.n, m: tc.m, k: tc.k, cells: make([]int, expect)}
	for i, t := range tokens {
		val, err := strconv.Atoi(t)
		if err != nil {
			return grid{}, false, fmt.Errorf("invalid integer %q: %v", t, err)
		}
		if val < 0 || val > 200000 {
			return grid{}, false, fmt.Errorf("value out of range: %d", val)
		}
		g.cells[i] = val
	}
	return g, false, nil
}

func validateGrid(tc testCase, g grid) error {
	if g.n != tc.n || g.m != tc.m || g.k != tc.k {
		return fmt.Errorf("grid dimensions mismatch")
	}
	idx := func(x, y, z int) int {
		return (x*g.m+y)*g.k + z
	}

	firstNonZero := func(vals []int) int {
		for _, v := range vals {
			if v != 0 {
				return v
			}
		}
		return 0
	}

	for y := 0; y < tc.m; y++ {
		for z := 0; z < tc.k; z++ {
			line := make([]int, tc.n)
			for x := 0; x < tc.n; x++ {
				line[x] = g.cells[idx(x, y, z)]
			}
			if firstNonZero(line) != tc.s.l[y][z] {
				return fmt.Errorf("sensor (0,%d,%d) mismatch", y+1, z+1)
			}
			reverse := make([]int, tc.n)
			for i := 0; i < tc.n; i++ {
				reverse[i] = line[tc.n-1-i]
			}
			if firstNonZero(reverse) != tc.s.r[y][z] {
				return fmt.Errorf("sensor (n+1,%d,%d) mismatch", y+1, z+1)
			}
		}
	}

	for x := 0; x < tc.n; x++ {
		for z := 0; z < tc.k; z++ {
			line := make([]int, tc.m)
			for y := 0; y < tc.m; y++ {
				line[y] = g.cells[idx(x, y, z)]
			}
			if firstNonZero(line) != tc.s.f[x][z] {
				return fmt.Errorf("sensor (%d,0,%d) mismatch", x+1, z+1)
			}
			reverse := make([]int, tc.m)
			for i := 0; i < tc.m; i++ {
				reverse[i] = line[tc.m-1-i]
			}
			if firstNonZero(reverse) != tc.s.b[x][z] {
				return fmt.Errorf("sensor (%d,m+1,%d) mismatch", x+1, z+1)
			}
		}
	}

	for x := 0; x < tc.n; x++ {
		for y := 0; y < tc.m; y++ {
			line := make([]int, tc.k)
			for z := 0; z < tc.k; z++ {
				line[z] = g.cells[idx(x, y, z)]
			}
			if firstNonZero(line) != tc.s.d[x][y] {
				return fmt.Errorf("sensor (%d,%d,0) mismatch", x+1, y+1)
			}
			reverse := make([]int, tc.k)
			for i := 0; i < tc.k; i++ {
				reverse[i] = line[tc.k-1-i]
			}
			if firstNonZero(reverse) != tc.s.u[x][y] {
				return fmt.Errorf("sensor (%d,%d,k+1) mismatch", x+1, y+1)
			}
		}
	}
	return nil
}

func computeSensors(g grid) sensorsData {
	idx := func(x, y, z int) int {
		return (x*g.m+y)*g.k + z
	}
	s := sensorsData{
		l: make([][]int, g.m),
		r: make([][]int, g.m),
		f: make([][]int, g.n),
		b: make([][]int, g.n),
		d: make([][]int, g.n),
		u: make([][]int, g.n),
	}
	for y := 0; y < g.m; y++ {
		s.l[y] = make([]int, g.k)
		s.r[y] = make([]int, g.k)
		for z := 0; z < g.k; z++ {
			for x := 0; x < g.n; x++ {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.l[y][z] = val
					break
				}
			}
			for x := g.n - 1; x >= 0; x-- {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.r[y][z] = val
					break
				}
			}
		}
	}
	for x := 0; x < g.n; x++ {
		s.f[x] = make([]int, g.k)
		s.b[x] = make([]int, g.k)
		for z := 0; z < g.k; z++ {
			for y := 0; y < g.m; y++ {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.f[x][z] = val
					break
				}
			}
			for y := g.m - 1; y >= 0; y-- {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.b[x][z] = val
					break
				}
			}
		}
	}
	for x := 0; x < g.n; x++ {
		s.d[x] = make([]int, g.m)
		s.u[x] = make([]int, g.m)
		for y := 0; y < g.m; y++ {
			for z := 0; z < g.k; z++ {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.d[x][y] = val
					break
				}
			}
			for z := g.k - 1; z >= 0; z-- {
				val := g.cells[idx(x, y, z)]
				if val != 0 {
					s.u[x][y] = val
					break
				}
			}
		}
	}
	return s
}

func deterministicTests() []testCase {
	tests := make([]testCase, 0, 4)

	// Completely empty single cell.
	tests = append(tests, fromGrid(grid{
		n: 1, m: 1, k: 1,
		cells: []int{0},
	}))

	// Single block visible from all directions.
	tests = append(tests, fromGrid(grid{
		n: 1, m: 1, k: 1,
		cells: []int{1337},
	}))

	// Small 2x2x2 with mixed blocks.
	g := grid{n: 2, m: 2, k: 2, cells: make([]int, 8)}
	g.cells[idx3(g, 0, 0, 0)] = 5
	g.cells[idx3(g, 1, 1, 1)] = 7
	g.cells[idx3(g, 0, 1, 1)] = 3
	tests = append(tests, fromGrid(g))

	// Intentionally inconsistent sensors.
	bad := fromGrid(grid{
		n: 1, m: 1, k: 1,
		cells: []int{0},
	})
	bad.s.r[0][0] = 42
	tests = append(tests, bad)

	return tests
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	for len(tests) < cap(tests) {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		k := rng.Intn(4) + 1
		if n*m*k > 120 {
			continue
		}
		cells := make([]int, n*m*k)
		for i := 0; i < len(cells); i++ {
			if rng.Intn(100) < 45 {
				cells[i] = rng.Intn(12) + 1
			}
		}
		base := fromGrid(grid{n: n, m: m, k: k, cells: cells})
		tests = append(tests, base)

		// Occasionally add a contradictory variant to test impossible handling.
		if rng.Intn(3) == 0 {
			bad := cloneTest(base)
			if bad.s.l[0][0] == 0 {
				bad.s.r[0][0] = 1
			} else {
				bad.s.l[0][0] = 0
			}
			tests = append(tests, bad)
		}
	}
	return tests
}

func cloneTest(tc testCase) testCase {
	copyMat := func(src [][]int) [][]int {
		dst := make([][]int, len(src))
		for i := range src {
			dst[i] = append([]int(nil), src[i]...)
		}
		return dst
	}
	return testCase{
		n: tc.n, m: tc.m, k: tc.k,
		s: sensorsData{
			l: copyMat(tc.s.l),
			r: copyMat(tc.s.r),
			f: copyMat(tc.s.f),
			b: copyMat(tc.s.b),
			d: copyMat(tc.s.d),
			u: copyMat(tc.s.u),
		},
	}
}

func fromGrid(g grid) testCase {
	return testCase{
		n: g.n,
		m: g.m,
		k: g.k,
		s: computeSensors(g),
	}
}

func idx3(g grid, x, y, z int) int {
	return (x*g.m+y)*g.k + z
}

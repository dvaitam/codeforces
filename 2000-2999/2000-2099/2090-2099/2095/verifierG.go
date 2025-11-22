package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x int64
	y int64
}

type testCase struct {
	n   int
	k   int
	pts []point
}

const (
	maxCoord   = int64(1_000_000_000)
	maxTotalN  = 120000
	maxRandN   = 4000
	largeN     = 100000
	randSeed   = 2095
	tolerance  = 1e-6
	safetyEps  = 1e-12
	usageGuide = "usage: go run verifierG.go /path/to/candidate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usageGuide)
		return
	}
	cand := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Printf("failed to build oracle: %v\n", err)
		return
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		oracleOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		expected, err := parseSingleFloat(oracleOut)
		if err != nil {
			fmt.Printf("oracle output parse error on test %d: %v\noutput:\n%s", idx+1, err, oracleOut)
			return
		}

		candOut, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		got, err := parseSingleFloat(candOut)
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			return
		}

		if !closeEnough(expected, got) {
			fmt.Printf("Mismatch on test %d: expected %.15f, got %.15f\n", idx+1, expected, got)
			fmt.Println("Input used:")
			fmt.Print(input)
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2095G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2095G")
	cmd := exec.Command("go", "build", "-o", outPath, "2095G.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSingleFloat(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 token, got %d", len(fields))
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func closeEnough(expected, got float64) bool {
	diff := math.Abs(expected - got)
	den := math.Max(1.0, math.Abs(expected))
	return diff <= tolerance*den+safetyEps
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n * 32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(tc.k))
	sb.WriteByte('\n')
	for _, p := range tc.pts {
		sb.WriteString(strconv.FormatInt(p.x, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(p.y, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testCase {
	var tests []testCase
	totalN := 0
	rng := rand.New(rand.NewSource(randSeed))

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic edge and sample-like cases
	add(testCase{
		n:   1,
		k:   1,
		pts: []point{{0, 0}},
	})

	add(testCase{
		n: 2,
		k: 1,
		pts: []point{
			{0, 0},
			{5, 0},
		},
	})

	add(testCase{
		n: 3,
		k: 2,
		pts: []point{
			{0, 0},
			{100, 0},
			{2, 0},
		},
	})

	add(testCase{
		n: 5,
		k: 5,
		pts: []point{
			{-10, -10},
			{-5, -5},
			{0, 0},
			{5, 5},
			{10, 10},
		},
	})

	// Random collinear cases respecting bounds and uniqueness
	for totalN < maxTotalN && len(tests) < 150 {
		remain := maxTotalN - totalN
		maxN := maxRandN
		if remain < maxN {
			maxN = remain
		}
		n := rng.Intn(maxN-1) + 2 // at least 2 to keep direction meaningful
		k := rng.Intn(n) + 1

		tc := randomCollinearCase(n, k, rng)
		add(tc)
	}

	// One large stress case
	if totalN+largeN <= maxTotalN {
		tc := randomCollinearCase(largeN, largeN/2, rng)
		add(tc)
	}

	return tests
}

func randomCollinearCase(n, k int, rng *rand.Rand) testCase {
	for {
		dx := int64(rng.Intn(2001) - 1000)
		dy := int64(rng.Intn(2001) - 1000)
		if dx == 0 && dy == 0 {
			continue
		}
		baseX := int64(rng.Intn(1_800_001) - 900_000)
		baseY := int64(rng.Intn(1_800_001) - 900_000)

		posSet := make(map[int64]struct{}, n)
		for len(posSet) < n {
			p := int64(rng.Intn(2_000_001) - 1_000_000)
			posSet[p] = struct{}{}
		}
		positions := make([]int64, 0, n)
		for p := range posSet {
			positions = append(positions, p)
		}
		sort.Slice(positions, func(i, j int) bool { return positions[i] < positions[j] })

		pts := make([]point, n)
		ok := true
		for i, p := range positions {
			x := baseX + p*dx
			y := baseY + p*dy
			if abs64(x) > maxCoord || abs64(y) > maxCoord {
				ok = false
				break
			}
			pts[i] = point{x: x, y: y}
		}
		if !ok {
			continue
		}
		return testCase{n: n, k: k, pts: pts}
	}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

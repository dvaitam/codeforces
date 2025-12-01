package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "1090C.go"

type testCase struct {
	name  string
	input string
}

type inputData struct {
	n     int
	m     int
	total int
	boxes [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		data, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse internal test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s\nreference output:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		minMoves, err := parseMoveCount(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\nreference output:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if err := verifyCandidate(candOut, data, minMoves); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\ncandidate output:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "1090C-verifier-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseMoveCount(out string) (int, error) {
	reader := strings.NewReader(out)
	var moves int
	if _, err := fmt.Fscan(reader, &moves); err != nil {
		return 0, fmt.Errorf("failed to read move count: %v", err)
	}
	if moves < 0 {
		return 0, fmt.Errorf("negative move count %d", moves)
	}
	return moves, nil
}

func parseInput(input string) (inputData, error) {
	reader := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return inputData{}, err
	}
	boxes := make([][]int, n)
	total := 0
	for i := 0; i < n; i++ {
		var s int
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return inputData{}, err
		}
		boxes[i] = make([]int, s)
		for j := 0; j < s; j++ {
			if _, err := fmt.Fscan(reader, &boxes[i][j]); err != nil {
				return inputData{}, err
			}
		}
		total += s
	}
	return inputData{n: n, m: m, total: total, boxes: boxes}, nil
}

func verifyCandidate(output string, data inputData, minMoves int) error {
	reader := strings.NewReader(output)
	var declared int
	if _, err := fmt.Fscan(reader, &declared); err != nil {
		return fmt.Errorf("failed to read move count: %v", err)
	}
	if declared < 0 {
		return fmt.Errorf("move count must be non-negative, got %d", declared)
	}
	if declared != minMoves {
		return fmt.Errorf("move count mismatch: expected %d, got %d", minMoves, declared)
	}

	states := make([]map[int]struct{}, data.n)
	sizes := make([]int, data.n)
	for i := 0; i < data.n; i++ {
		state := make(map[int]struct{}, len(data.boxes[i]))
		for _, val := range data.boxes[i] {
			state[val] = struct{}{}
		}
		states[i] = state
		sizes[i] = len(data.boxes[i])
	}

	for move := 0; move < declared; move++ {
		var from, to, kind int
		if _, err := fmt.Fscan(reader, &from, &to, &kind); err != nil {
			return fmt.Errorf("failed to read move %d: %v", move+1, err)
		}
		if from < 1 || from > data.n || to < 1 || to > data.n {
			return fmt.Errorf("move %d references invalid box index", move+1)
		}
		if kind < 1 || kind > data.m {
			return fmt.Errorf("move %d references invalid present kind", move+1)
		}
		from--
		to--
		if _, ok := states[from][kind]; !ok {
			return fmt.Errorf("move %d: box %d does not contain present %d", move+1, from+1, kind)
		}
		if _, ok := states[to][kind]; ok {
			return fmt.Errorf("move %d: box %d already contains present %d", move+1, to+1, kind)
		}
		delete(states[from], kind)
		states[to][kind] = struct{}{}
		sizes[from]--
		sizes[to]++
	}

	base := 0
	extra := 0
	if data.n > 0 {
		base = data.total / data.n
		extra = data.total % data.n
	}

	countExtra := 0
	for i, size := range sizes {
		if size == base {
			continue
		}
		if extra > 0 && size == base+1 {
			countExtra++
			continue
		}
		high := base
		if extra > 0 {
			high = base + 1
		}
		return fmt.Errorf("box %d has size %d, expected %d or %d", i+1, size, base, high)
	}

	if countExtra != extra {
		return fmt.Errorf("expected %d boxes of size %d, found %d", extra, base+1, countExtra)
	}

	var leftover string
	if _, err := fmt.Fscan(reader, &leftover); err == nil {
		return fmt.Errorf("unexpected extra output after listed moves")
	}

	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1090003))
	var tests []testCase

	tests = append(tests, makeCase("sample", 3, 5, [][]int{
		{1, 2, 3, 4, 5},
		{1, 2},
		{3, 4},
	}))
	tests = append(tests, makeCase("balanced_already", 4, 8, [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
	}))
	tests = append(tests, makeCase("all_empty", 5, 5, [][]int{
		{}, {}, {}, {}, {},
	}))
	tests = append(tests, makeCase("single_box", 1, 6, [][]int{
		{1, 2, 3, 4, 5, 6},
	}))
	tests = append(tests, makeCase("skewed", 6, 20, [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{},
		{8},
		{9, 10},
		{11, 12, 13, 14},
		{15, 16, 17},
	}))
	tests = append(tests, makeCase("duplicate_pressure", 5, 12, [][]int{
		{1, 2, 3, 4},
		{2, 5, 6},
		{3, 7},
		{4, 8, 9},
		{10, 11, 12},
	}))

	for i := 0; i < 5; i++ {
		n := rng.Intn(6) + 2 // 2..7
		m := rng.Intn(8) + 3 // 3..10
		maxTotal := n * m
		total := rng.Intn(min(60, maxTotal) + 1)
		tests = append(tests, randomCase(fmt.Sprintf("random_small_%d", i+1), rng, n, m, total))
	}

	for i := 0; i < 4; i++ {
		n := rng.Intn(30) + 20 // 20..49
		m := rng.Intn(80) + 20 // 20..99
		maxTotal := n * m
		total := rng.Intn(min(800, maxTotal-1)) + 1
		tests = append(tests, randomCase(fmt.Sprintf("random_mid_%d", i+1), rng, n, m, total))
	}

	for i := 0; i < 3; i++ {
		n := rng.Intn(200) + 100 // 100..299
		m := rng.Intn(500) + 200 // 200..699
		maxTotal := n * m
		total := rng.Intn(min(20000, maxTotal-1)) + 1
		tests = append(tests, randomCase(fmt.Sprintf("random_large_%d", i+1), rng, n, m, total))
	}

	tests = append(tests, randomCase("wide_1k", rng, 1000, 2000, 200000))
	tests = append(tests, randomCase("dense_small_m", rng, 2000, 120, 150000))
	tests = append(tests, randomCase("huge_sparse", rng, 100000, 100000, 500000))

	return tests
}

func makeCase(name string, n, m int, boxes [][]int) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for _, box := range boxes {
		fmt.Fprintf(&b, "%d", len(box))
		for _, val := range box {
			fmt.Fprintf(&b, " %d", val)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(name string, rng *rand.Rand, n, m, total int) testCase {
	if n == 0 {
		return makeCase(name, n, m, nil)
	}
	if total > n*m {
		total = n * m
	}
	boxes := make([][]int, n)
	rem := total
	for i := 0; i < n; i++ {
		if rem == 0 {
			boxes[i] = nil
			continue
		}
		remaining := n - i
		lower64 := int64(rem) - int64(remaining-1)*int64(m)
		if lower64 < 0 {
			lower64 = 0
		}
		upper64 := int64(rem)
		if upper64 > int64(m) {
			upper64 = int64(m)
		}
		if lower64 > upper64 {
			lower64 = upper64
		}
		var size int
		if remaining == 1 {
			size = rem
		} else {
			lo := int(lower64)
			hi := int(upper64)
			if hi < lo {
				hi = lo
			}
			if hi == lo {
				size = lo
			} else {
				size = rng.Intn(hi-lo+1) + lo
			}
		}
		if size > m {
			size = m
		}
		if size > 0 {
			seed := rng.Intn(m) + i + 1
			boxes[i] = uniqueRange(size, seed, m)
		}
		rem -= size
	}
	return makeCase(name, n, m, boxes)
}

func uniqueRange(size, seed, m int) []int {
	if size == 0 {
		return nil
	}
	if m <= 0 {
		panic("m must be positive")
	}
	arr := make([]int, size)
	start := ((seed % m) + m) % m
	for i := 0; i < size; i++ {
		arr[i] = ((start + i) % m) + 1
	}
	return arr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

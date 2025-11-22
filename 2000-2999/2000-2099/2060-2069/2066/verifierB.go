package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const randomTests = 200
const maxTotalN = 5000

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1], "candidate2066B")
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, test := range tests {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("deterministic test %d failed: %v\ninput:\n%s", idx+1, err, test.input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		test := randomTest(rng)
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("random test %d failed: %v\ninput:\n%s", i+1, err, test.input)
			return
		}
		total++
	}

	for idx, test := range largeTests() {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("large test %d failed: %v\ninput length: %d bytes\n", idx+1, err, len(test.input))
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path, prefix string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2066B.go")
	bin := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2066B_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return bin, func() { os.Remove(bin) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	expectedOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	gotOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expected, err := parseOutput(expectedOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseOutput(gotOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}

	for i := 0; i < test.t; i++ {
		if expected[i] != got[i] {
			return fmt.Errorf("case %d: expected %d got %d", i+1, expected[i], got[i])
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func parseOutput(out string, t int) ([]int, error) {
	reader := strings.NewReader(out)
	res := make([]int, 0, t)
	for len(res) < t {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return nil, fmt.Errorf("need %d integers, got %d (%v)", t, len(res), err)
		}
		res = append(res, x)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d integers, output has extra data", t)
	}
	return res, nil
}

func deterministicTests() []testInput {
	first := buildInput([]caseSpec{
		{n: 5, arr: []int{4, 3, 2, 1, 0}},
		{n: 6, arr: []int{4, 3, 3, 2, 1, 0}},
		{n: 4, arr: []int{2, 0, 1, 2}},
		{n: 1, arr: []int{7}},
	})
	second := buildInput([]caseSpec{
		{n: 5, arr: []int{0, 0, 0, 0, 0}},
		{n: 7, arr: []int{5, 5, 1, 0, 2, 3, 4}},
		{n: 7, arr: []int{9, 8, 7, 6, 5, 4, 3}},
		{n: 10, arr: []int{0, 5, 5, 5, 5, 5, 0, 1, 2, 3}},
	})
	return []testInput{first, second}
}

type caseSpec struct {
	n   int
	arr []int
}

func buildInput(cases []caseSpec) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d\n", c.n)
		writeArray(&sb, c.arr)
	}
	return testInput{input: sb.String(), t: len(cases)}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(5) + 1
	cases := make([]caseSpec, 0, t)
	totalN := 0
	for len(cases) < t && totalN < maxTotalN {
		remain := maxTotalN - totalN
		nLimit := min(200, remain)
		if nLimit <= 0 {
			break
		}
		n := rng.Intn(nLimit) + 1
		cases = append(cases, caseSpec{
			n:   n,
			arr: randomArray(n, rng),
		})
		totalN += n
	}
	return buildInput(cases)
}

func randomArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	mode := rng.Intn(6)
	switch mode {
	case 0: // increasing small values with zero present
		for i := 0; i < n; i++ {
			arr[i] = i
		}
	case 1: // descending values
		for i := 0; i < n; i++ {
			arr[i] = n - i
		}
	case 2: // random small range
		limit := rng.Intn(n+3) + 1
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(limit)
		}
	case 3: // large numbers mixed with a few zeros
		for i := 0; i < n; i++ {
			if rng.Intn(5) == 0 {
				arr[i] = 0
			} else {
				arr[i] = rng.Intn(1_000_000_000)
			}
		}
	case 4: // zero heavy prefix then random tail
		zeros := rng.Intn(n)
		for i := 0; i < zeros; i++ {
			arr[i] = 0
		}
		for i := zeros; i < n; i++ {
			arr[i] = rng.Intn(n + 2)
		}
	default: // alternating pattern
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				arr[i] = 0
			} else {
				arr[i] = (i % 7) + 1
			}
		}
	}
	return arr
}

func largeTests() []testInput {
	large1 := buildInput([]caseSpec{
		{n: 200000, arr: increasing(200000)},
	})
	large2 := buildInput([]caseSpec{
		{n: 100000, arr: zeroHeavy(100000)},
		{n: 100000, arr: mixedLarge(100000)},
	})
	return []testInput{large1, large2}
}

func increasing(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	return arr
}

func zeroHeavy(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if i%3 == 0 {
			arr[i] = 0
		} else {
			arr[i] = (i % 5) + 1
		}
	}
	return arr
}

func mixedLarge(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if i%4 == 0 {
			arr[i] = 0
		} else if i%4 == 1 {
			arr[i] = rngBounded(i)
		} else {
			arr[i] = 1_000_000_000 - (i % 1000)
		}
	}
	return arr
}

func rngBounded(v int) int {
	limit := (v % 100) + 2
	return v % limit
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

package main

import (
	"bytes"
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

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1])
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

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidate2066A_%d", time.Now().UnixNano()))
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
	src := filepath.Join(dir, "2066A.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2066A_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	oracleOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	candOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expect, err := parseTypes(oracleOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseTypes(candOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}
	for i := 0; i < test.t; i++ {
		if expect[i] != got[i] {
			return fmt.Errorf("case %d: expected %c got %c", i+1, expect[i], got[i])
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseTypes(output string, expected int) ([]byte, error) {
	reader := strings.NewReader(output)
	res := make([]byte, 0, expected)
	for len(res) < expected {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return nil, fmt.Errorf("need %d tokens, got %d (%v)", expected, len(res), err)
		}
		if len(token) == 0 {
			return nil, fmt.Errorf("empty token")
		}
		ch := strings.ToUpper(token)[0]
		if ch != 'A' && ch != 'B' {
			return nil, fmt.Errorf("invalid token %q", token)
		}
		res = append(res, ch)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d tokens, output has extra data", expected)
	}
	return res, nil
}

func deterministicTests() []testInput {
	example := `2
3
2 2 3
1 3 1
A
5
5 1 4 2 3
3 3 2 4 1
B
`
	small := buildTestInput([]caseSpec{
		newCaseSpec(3, []int{1, 2, 3}, []int{2, 3, 1}, 'A'),
		newCaseSpec(4, []int{4, 4, 4, 4}, []int{1, 2, 3, 1}, 'B'),
	})
	return []testInput{
		{input: example, t: 2},
		small,
	}
}

type caseSpec struct {
	n   int
	x   []int
	y   []int
	typ byte
}

func newCaseSpec(n int, x, y []int, typ byte) caseSpec {
	return caseSpec{n: n, x: append([]int(nil), x...), y: append([]int(nil), y...), typ: typ}
}

func buildTestInput(cases []caseSpec) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d\n", c.n)
		writeArray(&sb, c.x)
		writeArray(&sb, c.y)
		sb.WriteByte(byte(c.typ))
		sb.WriteByte('\n')
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
	cases := make([]caseSpec, t)
	totalN := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 3
		if totalN+n > 60 {
			n = 3
		}
		x, y := randomPairs(rng, n)
		typ := byte('A')
		if rng.Intn(2) == 1 {
			typ = 'B'
		}
		cases[i] = caseSpec{n: n, x: x, y: y, typ: typ}
		totalN += n
	}
	return buildTestInput(cases)
}

func randomPairs(rng *rand.Rand, n int) ([]int, []int) {
	x := make([]int, n)
	y := make([]int, n)
	used := make(map[int]struct{}, n)
	for i := 0; i < n; i++ {
		for {
			xi := rng.Intn(n) + 1
			yi := rng.Intn(n) + 1
			if xi == yi {
				continue
			}
			key := xi*(n+1) + yi
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			x[i] = xi
			y[i] = yi
			break
		}
	}
	return x, y
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2120-2129/2126/2126F.go"

type testCase struct {
	name    string
	input   string
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if !equalSlices(refAns, candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sreference:\n%s\ncandidate:\n%s", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2126F-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2126F.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		buildSmall("line3", []int{1, 2, 1}, [][3]int{{1, 2, 3}, {2, 3, 4}}, [][2]int{{1, 2}, {2, 1}, {3, 2}}),
		buildSmall("star4", []int{1, 1, 2, 2}, [][3]int{{1, 2, 5}, {1, 3, 2}, {1, 4, 3}}, [][2]int{{1, 2}, {4, 1}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSmall(name string, colors []int, edges [][3]int64, queries [][2]int) testCase {
	n := len(colors)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(queries)))
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(c))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return testCase{name: name, input: sb.String(), answers: len(queries)}
}

func newTestCase(name, input string) testCase {
	cnt, err := countAnswers(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, answers: cnt}
}

func countAnswers(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, fmt.Errorf("failed to read t: %v", err)
	}
	if t <= 0 {
		return 0, fmt.Errorf("non-positive t: %d", t)
	}
	total := 0
	for ; t > 0; t-- {
		var n, q int
		if _, err := fmt.Fscan(reader, &n, &q); err != nil {
			return 0, fmt.Errorf("failed to read n q: %v", err)
		}
		for i := 0; i < n; i++ {
			var tmp int
			if _, err := fmt.Fscan(reader, &tmp); err != nil {
				return 0, fmt.Errorf("failed to read color %d: %v", i, err)
			}
		}
		for i := 0; i < n-1; i++ {
			var u, v int
			var w int64
			if _, err := fmt.Fscan(reader, &u, &v, &w); err != nil {
				return 0, fmt.Errorf("failed to read edge %d: %v", i, err)
			}
		}
		for i := 0; i < q; i++ {
			var v, x int
			if _, err := fmt.Fscan(reader, &v, &x); err != nil {
				return 0, fmt.Errorf("failed to read query %d: %v", i, err)
			}
		}
		total += q
	}
	return total, nil
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	totalAns := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 2
		q := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		colors := make([]int, n)
		for j := 0; j < n; j++ {
			colors[j] = rng.Intn(3) + 1
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(colors[j]))
		}
		sb.WriteByte('\n')
		edges := generateTree(rng, n)
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		for j := 0; j < q; j++ {
			v := rng.Intn(n) + 1
			newColor := rng.Intn(3) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", v, newColor))
		}
		totalAns += q
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		answers: totalAns,
	}
}

func generateTree(rng *rand.Rand, n int) [][3]int64 {
	edges := make([][3]int64, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		w := int64(rng.Intn(10) + 1)
		edges = append(edges, [3]int64{int64(u), int64(v), w})
	}
	return edges
}

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

// refSource points to the local reference implementation to avoid GOPATH resolution.
const refSource = "418E.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := compareOutputs(tc.input, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-418E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref418E.bin")
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

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func compareOutputs(input, refOut, candOut string) error {
	expectCount, err := countAnswerQueries(input)
	if err != nil {
		return err
	}
	refVals, err := parseAnswerList(refOut)
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}
	candVals, err := parseAnswerList(candOut)
	if err != nil {
		return fmt.Errorf("candidate output invalid: %v", err)
	}
	if len(refVals) != expectCount {
		return fmt.Errorf("reference returned %d answers but expected %d", len(refVals), expectCount)
	}
	if len(candVals) != expectCount {
		return fmt.Errorf("candidate returned %d answers but expected %d", len(candVals), expectCount)
	}
	for i := 0; i < expectCount; i++ {
		if candVals[i] != refVals[i] {
			return fmt.Errorf("answer %d mismatch: expected %d got %d", i+1, refVals[i], candVals[i])
		}
	}
	return nil
}

func countAnswerQueries(input string) (int, error) {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, fmt.Errorf("failed to read n: %v", err)
	}
	for i := 0; i < n; i++ {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			return 0, fmt.Errorf("failed to read initial row: %v", err)
		}
	}
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return 0, fmt.Errorf("failed to read m: %v", err)
	}
	answerCount := 0
	for i := 0; i < m; i++ {
		var typ int
		if _, err := fmt.Fscan(reader, &typ); err != nil {
			return 0, fmt.Errorf("failed to read query type: %v", err)
		}
		if typ == 1 {
			var v int64
			var p int
			fmt.Fscan(reader, &v, &p)
		} else if typ == 2 {
			var x int
			var y int
			fmt.Fscan(reader, &x, &y)
			answerCount++
		} else {
			return 0, fmt.Errorf("unknown query type %d", typ)
		}
	}
	return answerCount, nil
}

func parseAnswerList(out string) ([]int64, error) {
	if len(strings.TrimSpace(out)) == 0 {
		return []int64{}, nil
	}
	lines := strings.Fields(out)
	values := make([]int64, len(lines))
	for i, token := range lines {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		values[i] = val
	}
	return values, nil
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "single_query", input: formatCase([]int{5}, []op{
			{typ: 2, x: 1, y: 1},
		})},
		{name: "updates_and_queries", input: formatCase([]int{1, 1, 2}, []op{
			{typ: 2, x: 2, y: 2},
			{typ: 1, v: 3, pos: 1},
			{typ: 2, x: 1, y: 1},
			{typ: 2, x: 100000, y: 3},
		})},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	return tests
}

type op struct {
	typ int
	v   int
	pos int
	x   int
	y   int
}

func formatCase(initial []int, ops []op) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(initial))
	for i, v := range initial {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(ops))
	for _, op := range ops {
		if op.typ == 1 {
			fmt.Fprintf(&sb, "1 %d %d\n", op.v, op.pos)
		} else {
			fmt.Fprintf(&sb, "2 %d %d\n", op.x, op.y)
		}
	}
	return sb.String()
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(8) + 1
	initial := make([]int, n)
	for i := 0; i < n; i++ {
		initial[i] = rng.Intn(6) + 1
	}
	m := rng.Intn(25) + 1
	ops := make([]op, 0, m)
	for i := 0; i < m; i++ {
		if rng.Intn(3) == 0 {
			v := rng.Intn(10) + 1
			pos := rng.Intn(n) + 1
			ops = append(ops, op{typ: 1, v: v, pos: pos})
		} else {
			x := rng.Intn(100000) + 1
			y := rng.Intn(n) + 1
			ops = append(ops, op{typ: 2, x: x, y: y})
		}
	}
	name := fmt.Sprintf("random_%d", idx+1)
	return testCase{name: name, input: formatCase(initial, ops)}
}

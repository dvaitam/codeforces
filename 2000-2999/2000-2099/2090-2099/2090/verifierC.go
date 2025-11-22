package main

import (
	"bufio"
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

const (
	refSource2090C = "2000-2999/2000-2099/2090-2099/2090/2090C.go"
	maxTotalGuests = 50000
)

type testCase struct {
	t []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2090C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := serialize(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expect, err := parseOutput(refOut, len(tc.t))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, runErr := runProgram(candidate, input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, runErr, input)
			os.Exit(1)
		}
		got, parseErr := parseOutput(candOut, len(tc.t))
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, parseErr, input, candOut)
			os.Exit(1)
		}

		if err := compare(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2090C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func resolveSourcePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string, n int) ([][2]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([][2]int64, 0, n)
	for i := 0; i < n; i++ {
		var x, y int64
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return nil, fmt.Errorf("expected %d coordinate pairs, got %d (%v)", n, i, err)
		}
		res = append(res, [2]int64{x, y})
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return res, nil
}

func compare(expect, got [][2]int64) error {
	if len(expect) != len(got) {
		return fmt.Errorf("length mismatch: expected %d guests, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i][0] != got[i][0] || expect[i][1] != got[i][1] {
			return fmt.Errorf("guest %d: expected (%d, %d), got (%d, %d)", i+1, expect[i][0], expect[i][1], got[i][0], got[i][1])
		}
	}
	return nil
}

func serialize(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(tc.t)))
	sb.WriteByte('\n')
	for i, v := range tc.t {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := []testCase{
		{t: []int{0, 1, 1, 0, 0, 1}},
		{t: []int{1, 0, 0, 1, 1}},
		{t: []int{1}},
		{t: []int{0}},
		{t: []int{1, 1, 1, 1, 1}},
		{t: []int{0, 0, 0, 0, 0}},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	for _, tc := range tests {
		total += len(tc.t)
	}

	add := func(tc testCase) {
		tests = append(tests, tc)
		total += len(tc.t)
	}

	add(alternatingCase(20))
	add(alternatingCase(41))
	add(randomCase(rng, 100))
	add(randomCase(rng, 500))
	add(blocksCase(300, rng))

	for total < maxTotalGuests {
		remaining := maxTotalGuests - total
		n := rng.Intn(min(2000, remaining)) + 1
		add(randomCase(rng, n))
	}

	return tests
}

func randomCase(rng *rand.Rand, n int) testCase {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(2)
	}
	return testCase{t: arr}
}

func alternatingCase(n int) testCase {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i % 2
	}
	return testCase{t: arr}
}

func blocksCase(n int, rng *rand.Rand) testCase {
	arr := make([]int, n)
	cur := rng.Intn(2)
	lenLeft := 0
	for i := 0; i < n; i++ {
		if lenLeft == 0 {
			lenLeft = rng.Intn(7) + 1
			cur ^= 1
		}
		arr[i] = cur
		lenLeft--
	}
	return testCase{t: arr}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

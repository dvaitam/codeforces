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

const (
	refSource = "./316D1.go"
	mod       = 1000000007
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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
		expectRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}
		expect, err := parseAnswer(expectRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		got, err := parseAnswer(gotRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		if expect != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, expect, got, tc.input, expectRaw, gotRaw)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-316D1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref316D1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
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

func parseAnswer(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 || val >= mod {
		return 0, fmt.Errorf("answer %d out of range [0,%d)", val, mod)
	}
	return val, nil
}

func buildTests() []testCase {
	var tests []testCase
	for n := 1; n <= 5; n++ {
		total := 1
		for i := 0; i < n; i++ {
			total *= 3
		}
		for mask := 0; mask < total; mask++ {
			caps := make([]int, n)
			tmp := mask
			for i := 0; i < n; i++ {
				caps[i] = tmp % 3
				tmp /= 3
			}
			tests = append(tests, testCase{
				name:  fmt.Sprintf("enum_n%d_mask%d", n, mask),
				input: formatInput(caps),
			})
		}
	}
	manual := [][]int{
		{1},
		{2},
		{1, 1},
		{2, 2},
		{0, 1, 2},
		{2, 1, 2, 1},
		{0, 0, 0, 0, 0},
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
	}
	for idx, arr := range manual {
		tests = append(tests, testCase{
			name:  fmt.Sprintf("manual_%d", idx+1),
			input: formatInput(arr),
		})
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(3)
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d_n%d", i+1, n),
			input: formatInput(arr),
		})
	}
	return tests
}

func formatInput(arr []int) string {
	var sb strings.Builder
	sb.Grow(len(arr)*4 + 20)
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

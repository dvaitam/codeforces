package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2140-2149/2144/2144E1.go"

type testCase struct {
	n    int
	arr  []int
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch: expected %d got %d\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCase(tc))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2144E1.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, arr []int) {
		cpy := append([]int(nil), arr...)
		tests = append(tests, testCase{n: len(cpy), arr: cpy, name: name})
	}

	add("single", []int{5})
	add("two_inc", []int{1, 2})
	add("two_same", []int{3, 3})
	add("small_mix", []int{3, 5, 5, 7, 4, 6, 7, 2, 4})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	const maxTotalN = 20000
	for len(tests) < 180 && totalN < maxTotalN {
		n := rng.Intn(100) + 1
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
			if n <= 0 {
				break
			}
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Intn(5) == 0 {
				arr[i] = rng.Intn(5) + 1
			} else {
				arr[i] = rng.Intn(1_000_000_000) + 1
			}
		}
		add(fmt.Sprintf("random_%d", len(tests)), arr)
		totalN += n
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, s := range fields {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = v
	}
	return ans, nil
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

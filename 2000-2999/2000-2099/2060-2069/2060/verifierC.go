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

const refSource = "2000-2999/2000-2099/2060-2069/2060/2060C.go"

type testCase struct {
	n    int
	k    int
	arr  []int
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	outPath := "./ref_2060C.bin"
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
	add := func(name string, k int, arr []int) {
		n := len(arr)
		cpy := append([]int(nil), arr...)
		tests = append(tests, testCase{n: n, k: k, arr: cpy, name: name})
	}

	add("sample1", 4, []int{1, 2, 3, 2})
	add("sample2", 15, []int{1, 2, 3, 4, 5, 6, 7, 8})
	add("sample3", 11, []int{1, 1, 1, 1, 1, 1})
	add("sample4", 9, []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9, 3})

	// Simple edges
	add("all_same_no_pair", 5, []int{1, 1, 1, 1})
	add("all_same_yes_pair", 2, []int{1, 1, 1, 1})
	add("alternating_pairs", 6, []int{1, 5, 1, 5, 1, 5, 1, 5})
	add("large_values", 20, []int{10, 10, 10, 10, 10, 10})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxTotalN = 180000
	totalN := 56 // sum of handcrafted cases above
	for len(tests) < 200 && totalN+50 <= maxTotalN {
		n := rng.Intn(20) + 2
		if rng.Intn(4) == 0 {
			n = rng.Intn(200) + 2
		}
		if n%2 == 1 {
			n++
		}
		if totalN+n > maxTotalN {
			break
		}
		k := rng.Intn(2*n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		add(fmt.Sprintf("random_%d", len(tests)), k, arr)
		totalN += n
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
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

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, s := range fields {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = v
	}
	return ans, nil
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

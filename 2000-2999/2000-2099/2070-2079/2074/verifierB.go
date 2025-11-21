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

const refSourceB = "2000-2999/2000-2099/2070-2079/2074/2074B.go"

type testCaseB struct {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReferenceB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTestsB()
	input := buildInputB(tests)

	refOut, err := runProgramB(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputB(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgramB(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputB(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch: expected %d got %d\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCaseB(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceB() (string, error) {
	outPath := "./ref_2074B.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSourceB)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgramB(target, input string) (string, error) {
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

func buildTestsB() []testCaseB {
	var tests []testCaseB
	add := func(name string, arr []int) {
		cpy := append([]int(nil), arr...)
		tests = append(tests, testCaseB{n: len(cpy), arr: cpy, name: name})
	}

	add("single_value", []int{10})
	add("two_values", []int{1, 2})
	add("all_ones_small", []int{1, 1, 1, 1})
	add("all_max_small", []int{1000, 1000, 1000})
	add("mixed_example", []int{998, 244, 353})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	const maxTotalN = 180000
	for len(tests) < 150 && totalN < maxTotalN {
		n := rng.Intn(400) + 1 // 1..400
		if len(tests)%20 == 0 {
			n = rng.Intn(2000) + 200 // occasional bigger
		}
		if totalN+n > maxTotalN {
			break
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1000) + 1
		}
		add(fmt.Sprintf("random_%d", len(tests)), arr)
		totalN += n
	}
	return tests
}

func buildInputB(tests []testCaseB) string {
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

func parseOutputB(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, s := range fields {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		res[i] = v
	}
	return res, nil
}

func formatCaseB(tc testCaseB) string {
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

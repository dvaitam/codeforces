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

type testCase struct {
	n   int
	k   int
	arr []int
}

const maxTotalN = 400000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTestSuite(rng)
	input := serializeInput(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, expected[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d k=%d arr=%v\n", tests[i].n, tests[i].k, tests[i].arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2085F2.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2085F2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, tests int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != tests {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", tests, len(fields), out)
	}
	res := make([]int64, tests)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
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
	return stdout.String(), nil
}

func buildTestSuite(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	total := totalLength(tests)

	for total < maxTotalN {
		remaining := maxTotalN - total
		if remaining < 2 {
			break
		}
		n := randomN(rng, remaining)
		if n < 2 {
			break
		}
		maxK := n
		if maxK > 2000 {
			maxK = 2000
		}
		if rng.Intn(5) == 0 && n >= 5000 {
			maxK = n
		}
		k := rng.Intn(maxK-1) + 2
		arr := generateArray(rng, n, k)
		tests = append(tests, testCase{n: n, k: k, arr: arr})
		total += n
	}
	return tests
}

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	return total
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 3, k: 2, arr: []int{1, 2, 1}},
		{n: 7, k: 3, arr: []int{2, 1, 1, 3, 1, 1, 2}},
		{n: 6, k: 3, arr: []int{1, 1, 2, 2, 2, 3}},
		{n: 6, k: 3, arr: []int{1, 2, 2, 2, 2, 3}},
		{n: 10, k: 5, arr: []int{5, 1, 3, 1, 1, 2, 2, 4, 1, 3}},
		{n: 9, k: 4, arr: []int{1, 2, 3, 3, 3, 3, 3, 2, 4}},
		{n: 5, k: 5, arr: []int{1, 2, 3, 4, 5}},
		{n: 8, k: 2, arr: []int{1, 1, 1, 1, 2, 2, 2, 2}},
	}
}

func randomN(rng *rand.Rand, remaining int) int {
	maxBlock := 8000
	if remaining < maxBlock {
		maxBlock = remaining
	}
	n := rng.Intn(maxBlock-1) + 2
	if n > remaining {
		n = remaining
	}
	return n
}

func generateArray(rng *rand.Rand, n, k int) []int {
	arr := make([]int, n)
	order := rng.Perm(k)
	for i := 0; i < k; i++ {
		arr[i] = order[i] + 1
	}
	for i := k; i < n; i++ {
		arr[i] = rng.Intn(k) + 1
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
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

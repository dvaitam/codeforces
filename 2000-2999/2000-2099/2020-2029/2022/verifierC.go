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
	n int
	s [2]string
}

const maxTotalN = 100000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
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

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nn=%d\nrow1=%s\nrow2=%s\n", i+1, expected[i], got[i], tests[i].n, tests[i].s[0], tests[i].s[1])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2022C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2022C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, cases int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != cases {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", cases, len(fields), out)
	}
	res := make([]int64, cases)
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

func buildTests() []testCase {
	tests := deterministicTests()
	total := totalLength(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < maxTotalN {
		remaining := maxTotalN - total
		n := rng.Intn(min(1000, remaining)) + 3
		n = (n / 3) * 3
		if n < 3 {
			n = 3
		}
		var rows [2]string
		for r := 0; r < 2; r++ {
			var sb strings.Builder
			for i := 0; i < n; i++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('A')
				} else {
					sb.WriteByte('J')
				}
			}
			rows[r] = sb.String()
		}
		tests = append(tests, testCase{n: n, s: rows})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 3, s: [2]string{"AAA", "JJJ"}},
		{n: 3, s: [2]string{"JJJ", "AAA"}},
		{n: 6, s: [2]string{"JAJAJJ", "JJAJAJ"}},
		{n: 6, s: [2]string{"AJJJAA", "JAJJAA"}},
	}
}

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	return total
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s[0])
		sb.WriteByte('\n')
		sb.WriteString(tc.s[1])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

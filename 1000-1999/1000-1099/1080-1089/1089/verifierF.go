package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesF.txt.
const testcasesFData = `6 720
4 24
4 24
8 40320
8 40320
9 362880
3 6
5 120
4 24
4 24
6 720
3 6
1 1
0 1
3 6
6 720
7 5040
8 40320
3 6
3 6
1 1
9 362880
6 720
1 1
0 1
2 2
5 120
0 1
6 720
3 6
10 3628800
1 1
9 362880
4 24
9 362880
9 362880
4 24
3 6
10 3628800
1 1
3 6
3 6
1 1
2 2
3 6
5 120
3 6
1 1
8 40320
2 2
7 5040
0 1
7 5040
2 2
4 24
9 362880
4 24
5 120
8 40320
2 2
9 362880
5 120
9 362880
3 6
2 2
4 24
3 6
6 720
2 2
4 24
7 5040
2 2
8 40320
4 24
4 24
1 1
5 120
10 3628800
3 6
3 6
1 1
7 5040
8 40320
4 24
3 6
5 120
2 2
3 6
4 24
2 2
7 5040
1 1
0 1
2 2
0 1
3 6
1 1
9 362880
3 6
0 1`

type testCase struct {
	n int
}

// solve mirrors the logic from 1089F.go for a single test case.
func solve(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	if len(tokens)%2 != 0 {
		return nil, fmt.Errorf("embedded testcases should have pairs of n and expected values")
	}
	cases := make([]testCase, 0, len(tokens)/2)
	for i := 0; i < len(tokens); i += 2 {
		n, err := strconv.Atoi(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("invalid n at token %d: %w", i, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d\n", tc.n)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.Itoa(solve(tc.n))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesFData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

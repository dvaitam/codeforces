package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// testCase represents one input with potentially multiple testcases
// t testcases with dimensions in ns and ms.
type testCase struct {
	input string
	ns    []int
	ms    []int
}

// randomCase generates a random test case with up to 3 individual boards of size up to 10x10.
func randomCase(rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	ns := make([]int, t)
	ms := make([]int, t)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(9) + 2 // 2..10
		m := rng.Intn(9) + 2
		ns[i] = n
		ms[i] = m
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	}
	return testCase{input: sb.String(), ns: ns, ms: ms}
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// checkOutput validates candidate output for the provided test case.
func checkOutput(tc testCase, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	idx := 0
	for caseNum := 0; caseNum < len(tc.ns); caseNum++ {
		n, m := tc.ns[caseNum], tc.ms[caseNum]
		if idx+n > len(lines) {
			return fmt.Errorf("case %d: expected %d lines, got %d", caseNum+1, n, len(lines)-idx)
		}
		board := make([][]byte, n)
		for i := 0; i < n; i++ {
			line := strings.TrimSpace(lines[idx+i])
			if len(line) != m {
				return fmt.Errorf("case %d line %d: expected length %d", caseNum+1, i+1, m)
			}
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				c := line[j]
				if c != 'B' && c != 'W' {
					return fmt.Errorf("case %d line %d: invalid char", caseNum+1, i+1)
				}
				row[j] = c
			}
			board[i] = row
		}
		idx += n
		// compute B and W
		B, W := 0, 0
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				here := board[i][j]
				for _, d := range dirs {
					x, y := i+d[0], j+d[1]
					if x >= 0 && x < n && y >= 0 && y < m {
						if here == 'B' && board[x][y] == 'W' {
							B++
							break
						}
						if here == 'W' && board[x][y] == 'B' {
							W++
							break
						}
					}
				}
			}
		}
		if B != W+1 {
			return fmt.Errorf("case %d: condition failed B=%d W=%d", caseNum+1, B, W)
		}
	}
	for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
		idx++
	}
	if idx != len(lines) {
		return fmt.Errorf("extra output lines")
	}
	return nil
}

func runCase(bin string, tc testCase) error {
	out, err := run(bin, tc.input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	return checkOutput(tc, out)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic simple cases
	cases := []testCase{
		{input: "1\n2 2\n", ns: []int{2}, ms: []int{2}},
		{input: "1\n3 3\n", ns: []int{3}, ms: []int{3}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

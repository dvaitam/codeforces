package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct{ input string }

func solveCase(in string) string {
	fields := strings.Fields(in)
	if len(fields) < 2 {
		return "0"
	}
	var n, m int
	fmt.Sscan(fields[0], &n)
	fmt.Sscan(fields[1], &m)
	lines := strings.Split(strings.TrimSpace(in), "\n")
	grid := lines[1 : 1+n]

	// dp[i][j] holds the number of consecutive 1's ending at column j in row i
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				if j == 0 {
					dp[i][j] = 1
				} else {
					dp[i][j] = dp[i][j-1] + 1
				}
			}
		}
	}

	maxArea := 0
	colHeights := make([]int, n)
	for j := 0; j < m; j++ {
		// gather column j widths and sort descending
		for i := 0; i < n; i++ {
			colHeights[i] = dp[i][j]
		}
		sort.Slice(colHeights, func(a, b int) bool { return colHeights[a] > colHeights[b] })
		for i := 0; i < n; i++ {
			area := colHeights[i] * (i + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return fmt.Sprint(maxArea)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.input)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	m := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		sb.Write(row)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 105)
	cases = append(cases, testCase{input: "1 1\n1\n"}, testCase{input: "2 2\n01\n10\n"})
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

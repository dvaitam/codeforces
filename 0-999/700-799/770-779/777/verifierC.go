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

type testCase struct {
	n, m    int
	grid    [][]int
	k       int
	queries [][2]int
}

func solve(tc testCase) []string {
	n, m := tc.n, tc.m
	a := tc.grid
	up := make([]int, m)
	rowMin := make([]int, n)
	for i := 0; i < n; i++ {
		minVal := i
		if i == 0 {
			for j := 0; j < m; j++ {
				up[j] = 0
			}
			minVal = 0
		} else {
			for j := 0; j < m; j++ {
				if a[i][j] < a[i-1][j] {
					up[j] = i
				}
				if up[j] < minVal {
					minVal = up[j]
				}
			}
		}
		if i == 0 {
			minVal = 0
		}
		rowMin[i] = minVal
	}
	res := make([]string, len(tc.queries))
	for idx, q := range tc.queries {
		l := q[0] - 1
		r := q[1] - 1
		if rowMin[r] <= l {
			res[idx] = "Yes"
		} else {
			res[idx] = "No"
		}
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", tc.k))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return sb.String()
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 110)
	for i := 0; i < 110; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			grid[r] = make([]int, m)
			for c := 0; c < m; c++ {
				grid[r][c] = rng.Intn(20)
			}
		}
		k := rng.Intn(n) + 1
		queries := make([][2]int, k)
		for j := 0; j < k; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[j] = [2]int{l, r}
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid, k: k, queries: queries})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expRes := solve(tc)
		exp := strings.Join(expRes, "\n")
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

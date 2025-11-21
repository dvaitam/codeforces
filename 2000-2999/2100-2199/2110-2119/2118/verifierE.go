package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	m int
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	total := 0
	for len(cases) < 50 && total < 2000 {
		n := rng.Intn(9)*2 + 1 // odd between 1 and 17
		m := rng.Intn(9)*2 + 1
		cases = append(cases, testCase{n: n, m: m})
		total += n * m
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, m: 1})
	}
	return cases
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

type cell struct {
	r, c int
}

func validateCase(tc testCase, coords []cell) error {
	n, m := tc.n, tc.m
	if len(coords) != n*m {
		return fmt.Errorf("expected %d coordinates got %d", n*m, len(coords))
	}
	colored := make([][]bool, n)
	for i := range colored {
		colored[i] = make([]bool, m)
	}
	penalty := make([][]int, n)
	for i := range penalty {
		penalty[i] = make([]int, m)
	}
	coloredList := make([]cell, 0, n*m)
	for step, p := range coords {
		r, c := p.r, p.c
		if r < 1 || r > n || c < 1 || c > m {
			return fmt.Errorf("step %d: coordinate out of bounds (%d,%d)", step+1, r, c)
		}
		r--
		c--
		if colored[r][c] {
			return fmt.Errorf("step %d: cell (%d,%d) already colored", step+1, r+1, c+1)
		}
		colored[r][c] = true
		coloredList = append(coloredList, cell{r, c})

		// find farthest colored cells from current
		maxCheb := -1
		maxMan := -1
		for _, q := range coloredList {
			dx := int(math.Abs(float64(q.r - r)))
			dy := int(math.Abs(float64(q.c - c)))
			cheb := dx
			if dy > cheb {
				cheb = dy
			}
			man := dx + dy
			if cheb > maxCheb || (cheb == maxCheb && man > maxMan) {
				maxCheb = cheb
				maxMan = man
			}
		}
		for _, q := range coloredList {
			dx := int(math.Abs(float64(q.r - r)))
			dy := int(math.Abs(float64(q.c - c)))
			cheb := dx
			if dy > cheb {
				cheb = dy
			}
			man := dx + dy
			if cheb == maxCheb && man == maxMan {
				penalty[q.r][q.c]++
				if penalty[q.r][q.c] > 3 {
					return fmt.Errorf("penalty exceeded at (%d,%d): %d", q.r+1, q.c+1, penalty[q.r][q.c])
				}
			}
		}
	}
	// ensure all cells colored
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if !colored[i][j] {
				return fmt.Errorf("cell (%d,%d) not colored", i+1, j+1)
			}
		}
	}
	return nil
}

func parseOutput(out string, cases []testCase) ([]cell, error) {
	tokens := strings.Fields(out)
	cur := 0
	var all []cell
	for _, tc := range cases {
		need := tc.n * tc.m * 2
		if cur+need > len(tokens) {
			return nil, fmt.Errorf("not enough tokens for a test case")
		}
		for i := 0; i < tc.n*tc.m; i++ {
			r, err1 := strconv.Atoi(tokens[cur])
			c, err2 := strconv.Atoi(tokens[cur+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("non-integer coordinate")
			}
			all = append(all, cell{r, c})
			cur += 2
		}
	}
	if cur != len(tokens) {
		return nil, fmt.Errorf("extra output tokens found")
	}
	return all, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/2118E_binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	cases := genCases()
	input := buildInput(cases)
	out, err := runCandidate(cand, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(out)
	cur := 0
	for idx, tc := range cases {
		need := tc.n * tc.m * 2
		if cur+need > len(tokens) {
			fmt.Fprintf(os.Stderr, "test %d: not enough coordinates\n", idx+1)
			os.Exit(1)
		}
		var coords []cell
		for i := 0; i < tc.n*tc.m; i++ {
			r, err1 := strconv.Atoi(tokens[cur])
			c, err2 := strconv.Atoi(tokens[cur+1])
			if err1 != nil || err2 != nil {
				fmt.Fprintf(os.Stderr, "test %d: invalid integer\n", idx+1)
				os.Exit(1)
			}
			coords = append(coords, cell{r: r, c: c})
			cur += 2
		}
		if err := validateCase(tc, coords); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v (n=%d m=%d)\n", idx+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	if cur != len(tokens) {
		fmt.Fprintf(os.Stderr, "extra output tokens detected\n")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

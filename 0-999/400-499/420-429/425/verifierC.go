package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseC struct {
	n, m, s, e int
	a          []int
	b          []int
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(8) + 2 //2..9
		m := rand.Intn(8) + 2
		s := rand.Intn(15) + 5 //5..19
		e := rand.Intn(3) + 1  //1..3
		a := make([]int, n)
		b := make([]int, m)
		for j := range a {
			a[j] = rand.Intn(5) + 1
		}
		for j := range b {
			b[j] = rand.Intn(5) + 1
		}
		tests[i] = testCaseC{n, m, s, e, a, b}
	}
	return tests
}

func solveC(tc testCaseC) int {
	n, m, s, e := tc.n, tc.m, tc.s, tc.e
	a, b := tc.a, tc.b
	maxVal := 0
	for _, v := range b {
		if v > maxVal {
			maxVal = v
		}
	}
	posList := make([][]int, maxVal+1)
	for i, v := range b {
		posList[v] = append(posList[v], i+1)
	}
	maxK := s / e
	if maxK > n {
		maxK = n
	}
	if maxK > m {
		maxK = m
	}
	INFJ := m + 1
	dp := make([]int, maxK+2)
	for i := range dp {
		dp[i] = INFJ
	}
	dp[0] = 0
	INF_SUM := n + m + 5
	minSum := make([]int, maxK+2)
	for i := range minSum {
		minSum[i] = INF_SUM
	}
	curMax := 0
	for i := 1; i <= n; i++ {
		v := a[i-1]
		if v <= maxVal {
			lst := posList[v]
			if len(lst) == 0 {
				continue
			}
			ub := curMax
			if ub > maxK-1 {
				ub = maxK - 1
			}
			for t := ub; t >= 0; t-- {
				prevJ := dp[t]
				if prevJ >= INFJ {
					continue
				}
				idx := sort.Search(len(lst), func(k int) bool { return lst[k] > prevJ })
				if idx < len(lst) {
					j := lst[idx]
					if j < dp[t+1] {
						dp[t+1] = j
					}
					sum := i + j
					if sum < minSum[t+1] {
						minSum[t+1] = sum
					}
					if t+1 > curMax {
						curMax = t + 1
					}
				}
			}
		}
	}
	ans := 0
	for t := 1; t <= curMax; t++ {
		if t*e > s {
			break
		}
		if minSum[t] <= s-t*e {
			ans = t
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d %d\n", tc.n, tc.m, tc.s, tc.e)
		for j, v := range tc.a {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for j, v := range tc.b {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expect := solveC(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

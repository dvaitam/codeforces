package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt to remove external dependency.
const testcasesRaw = `100
1 3
1
1 3
2 4
2
1 1
2 2
1
1 4
3 3
1
1 1
3
1 1
2 2
3 3
1
1 2
1 1
1
1 1
1 2
2
1 1
2 2
1 5
3
1 3
4 4
5 5
1 2
2
1 1
2 2
3 3
1
1 2
3
1 1
2 2
3 3
2
1 2
3 3
2 2
1
1 2
2
1 1
2 2
1 4
3
1 2
3 3
4 4
3 2
1
1 2
1
1 2
2
1 1
2 2
2 5
3
1 3
4 4
5 5
1
1 1
3 2
2
1 1
2 2
2
1 1
2 2
2
1 1
2 2
2 1
1
1 1
1
1 1
2 3
1
1 3
2
1 1
2 3
3 2
2
1 1
2 2
2
1 1
2 2
2
1 1
2 2
1 4
2
1 2
3 4
2 5
1
1 4
2
1 4
5 5
1 1
1
1 1
3 2
1
1 1
2
1 1
2 2
1
1 1
3 3
2
1 1
2 3
1
1 1
2
1 2
3 3
1 2
2
1 1
2 2
3 4
2
1 2
3 3
1
1 1
3
1 2
3 3
4 4
1 4
3
1 1
2 3
4 4
2 5
1
1 5
4
1 1
2 3
4 4
5 5
2 2
1
1 2
2
1 1
2 2
2 5
2
1 3
4 5
4
1 2
3 3
4 4
5 5
1 4
2
1 2
3 3
1 4
4
1 1
2 2
3 3
4 4
1 2
1
1 2
1 4
4
1 1
2 2
3 3
4 4
3 5
3
1 3
4 4
5 5
1
1 1
1
1 5
1 4
3
1 1
2 2
3 4
1 5
1
1 3
2 1
1
1 1
1
1 1
2 4
2
1 2
3 4
1
1 1
1 5
3
1 3
4 4
5 5
3 4
1
1 2
3
1 2
3 3
4 4
2
1 2
3 3
1 2
2
1 1
2 2
2 5
1
1 2
1
1 3
1 5
1
1 4
3 5
4
1 1
2 2
3 4
5 5
3
1 1
2 2
3 5
1
1 1
3 3
3
1 1
2 2
3 3
3
1 1
2 2
3 3
2
1 2
3 3
1 4
1
1 2
3 5
2
1 3
4 5
2
1 3
4 4
1
1 4
3 2
2
1 1
2 2
1
1 2
1
1 1
3 1
1
1 1
1
1 1
1
1 1
1 2
2
1 1
2 2
3 4
2
1 2
3 3
3
1 2
3 3
4 4
4
1 1
2 2
3 3
4 4
2 4
1
1 3
2
1 3
4 4
3 5
5
1 1
2 2
3 3
4 4
5 5
5
1 1
2 2
3 3
4 4
5 5
3
1 2
3 3
4 5
3 1
1
1 1
1
1 1
1
1 1
1 4
3
1 1
2 2
3 3
1 2
1
1 2
1 1
1
1 1
3 3
3
1 1
2 2
3 3
1
1 3
3
1 1
2 2
3 3
3 5
2
1 4
5 5
2
1 3
4 5
5
1 1
2 2
3 3
4 4
5 5
3 4
3
1 2
3 3
4 4
4
1 1
2 2
3 3
4 4
1
1 3
1 1
1
1 1
3 1
1
1 1
1
1 1
1
1 1
3 4
3
1 1
2 2
3 3
4
1 1
2 2
3 3
4 4
2
1 2
3 4
2 3
1
1 1
2
1 1
2 2
3 3
3
1 1
2 2
3 3
1
1 3
1
1 2
3 5
5
1 1
2 2
3 3
4 4
5 5
5
1 1
2 2
3 3
4 4
5 5
2
1 3
4 4
3 1
1
1 1
1
1 1
1
1 1
3 3
2
1 2
3 3
1
1 1
2
1 2
3 3
3 5
3
1 3
4 4
5 5
4
1 2
3 3
4 4
5 5
4
1 2
3 3
4 4
5 5
3 4
3
1 1
2 3
4 4
1
1 1
4
1 1
2 2
3 3
4 4
1 2
1
1 1
2 2
2
1 1
2 2
1
1 2
3 2
1
1 1
1
1 1
1
1 2
2 5
5
1 1
2 2
3 3
4 4
5 5
2
1 3
4 4
3 3
3
1 1
2 2
3 3
1
1 1
1
1 1
1 1
1
1 1
1 2
1
1 2
2 5
5
1 1
2 2
3 3
4 4
5 5
1
1 1
1 2
2
1 1
2 2
1 3
3
1 1
2 2
3 3
2 5
1
1 2
2
1 3
4 4
3 2
1
1 1
1
1 1
2
1 1
2 2
1 3
3
1 1
2 2
3 3
2 5
4
1 2
3 3
4 4
5 5
2
1 4
5 5
3 4
2
1 1
2 3
3
1 1
2 2
3 4
2
1 3
4 4
1 5
2
1 3
4 4
3 1
1
1 1
1
1 1
1
1 1
1 5
3
1 1
2 3
4 5
3 3
1
1 3
2
1 2
3 3
3
1 1
2 2
3 3
2 5
2
1 4
5 5
1
1 1
3 5
2
1 2
3 3
5
1 1
2 2
3 3
4 4
5 5
5
1 1
2 2
3 3
4 4
5 5
2 3
1
1 3
2
1 2
3 3
1 3
2
1 1
2 2
2 5
5
1 1
2 2
3 3
4 4
5 5
4
1 1
2 3
4 4
5 5
1 2
1
1 1
1 1
1
1 1
1 1
1
1 1
3 2
1
1 1
1
1 1
2
1 1
2 2
3 3
3
1 1
2 2
3 3
3
1 1
2 2
3 3
1
1 1
2 5
4
1 2
3 3
4 4
5 5
5
1 1
2 2
3 3
4 4
5 5
3 5
5
1 1
2 2
3 3
4 4
5 5
2
1 3
4 4
3
1 2
3 4
5 5
3 4
3
1 2
3 3
4 4
3
1 1
2 2
3 4
2
1 1
2 3`

type interval struct{ l, r int }

type testCase struct {
	n    int
	m    int
	rows [][]interval
}

// referenceSolution embeds the algorithm from 1372E.go.
func referenceSolution(n, m int, rows [][]interval) int {
	intervals := make([][2]int, 0)
	for _, row := range rows {
		for _, iv := range row {
			intervals = append(intervals, [2]int{iv.l, iv.r})
		}
	}
	m1 := m + 2
	arr := make([][][]int, m1)
	pre := make([][][]int, m1)
	for k := 0; k <= m; k++ {
		arr[k] = make([][]int, m1)
		pre[k] = make([][]int, m1)
		for i := 0; i <= m; i++ {
			arr[k][i] = make([]int, m1)
			pre[k][i] = make([]int, m1)
		}
	}
	for _, it := range intervals {
		L, R := it[0], it[1]
		for k := L; k <= R; k++ {
			arr[k][L][R]++
		}
	}
	for k := 1; k <= m; k++ {
		for l := m; l >= 1; l-- {
			for r := l; r <= m; r++ {
				v := arr[k][l][r]
				if l+1 <= m {
					v += pre[k][l+1][r]
				}
				if r-1 >= 1 {
					v += pre[k][l][r-1]
				}
				if l+1 <= m && r-1 >= 1 {
					v -= pre[k][l+1][r-1]
				}
				pre[k][l][r] = v
			}
		}
	}
	dp := make([][]int, m1)
	for i := range dp {
		dp[i] = make([]int, m1)
	}
	for length := 1; length <= m; length++ {
		for l := 1; l+length-1 <= m; l++ {
			r := l + length - 1
			best := 0
			for k := l; k <= r; k++ {
				cnt := pre[k][l][r]
				v := cnt * cnt
				if k > l {
					v += dp[l][k-1]
				}
				if k < r {
					v += dp[k+1][r]
				}
				if v > best {
					best = v
				}
			}
			dp[l][r] = best
		}
	}
	return dp[1][m]
}

func parseTestcases() []testCase {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		panic("no testcases")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			panic("missing n")
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			panic("missing m")
		}
		m, _ := strconv.Atoi(scan.Text())
		rows := make([][]interval, n)
		for r := 0; r < n; r++ {
			if !scan.Scan() {
				panic("missing row count")
			}
			k, _ := strconv.Atoi(scan.Text())
			row := make([]interval, k)
			for j := 0; j < k; j++ {
				if !scan.Scan() {
					panic("missing l")
				}
				l, _ := strconv.Atoi(scan.Text())
				if !scan.Scan() {
					panic("missing r")
				}
				rv, _ := strconv.Atoi(scan.Text())
				row[j] = interval{l: l, r: rv}
			}
			rows[r] = row
		}
		tests = append(tests, testCase{n: n, m: m, rows: rows})
	}
	return tests
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases()

	for idx, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.rows {
			fmt.Fprintf(&input, "%d\n", len(row))
			for _, iv := range row {
				fmt.Fprintf(&input, "%d %d\n", iv.l, iv.r)
			}
		}
		expect := referenceSolution(tc.n, tc.m, tc.rows)
		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Printf("case %d: invalid output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

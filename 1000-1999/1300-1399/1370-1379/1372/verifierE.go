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

func runCandidate(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type interval struct{ l, r int }

func solveCase(n, m int, rows [][]interval) int {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		rows := make([][]interval, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			k, _ := strconv.Atoi(scan.Text())
			row := make([]interval, k)
			for j := 0; j < k; j++ {
				scan.Scan()
				l, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				r, _ := strconv.Atoi(scan.Text())
				row[j] = interval{l: l, r: r}
			}
			rows[i] = row
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, row := range rows {
			fmt.Fprintf(&input, "%d\n", len(row))
			for _, iv := range row {
				fmt.Fprintf(&input, "%d %d\n", iv.l, iv.r)
			}
		}
		expect := solveCase(n, m, rows)
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Printf("case %d: invalid output\n", caseIdx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", caseIdx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

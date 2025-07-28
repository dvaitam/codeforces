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

func solveExpected(n int, tags []int, scores []int64) int64 {
	dp := make([]int64, n+1)
	var res int64
	for m := 1; m <= n; m++ {
		for j := m - 1; j >= 1; j-- {
			if tags[m] == tags[j] {
				continue
			}
			diff := scores[m] - scores[j]
			if diff < 0 {
				diff = -diff
			}
			oldJ := dp[j]
			if dp[m]+diff > dp[j] {
				dp[j] = dp[m] + diff
			}
			if oldJ+diff > dp[m] {
				dp[m] = oldJ + diff
			}
		}
		if dp[m] > res {
			res = dp[m]
		}
	}
	for i := 1; i <= n; i++ {
		if dp[i] > res {
			res = dp[i]
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			fmt.Printf("case %d invalid format\n", caseNum)
			os.Exit(1)
		}
		left := strings.Fields(strings.TrimSpace(parts[0]))
		right := strings.Fields(strings.TrimSpace(parts[1]))
		if len(left) < 1 {
			fmt.Printf("case %d invalid left part\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(left[0])
		if len(left)-1 != n || len(right) != n {
			fmt.Printf("case %d length mismatch\n", caseNum)
			os.Exit(1)
		}
		tags := make([]int, n+1)
		scores := make([]int64, n+1)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(left[1+i])
			tags[i+1] = v
		}
		for i := 0; i < n; i++ {
			vv, _ := strconv.Atoi(right[i])
			scores[i+1] = int64(vv)
		}
		exp := solveExpected(n, tags, scores)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(tags[i]))
		}
		input.WriteByte('\n')
		for i := 1; i <= n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(scores[i], 10))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		resultStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(resultStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: invalid output %q\n", caseNum, resultStr)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}

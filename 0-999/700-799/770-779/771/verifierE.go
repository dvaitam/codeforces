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

func solveCase(a, b []int64) int {
	n := len(a)
	prefix1, prefix2, prefixAll := int64(0), int64(0), int64(0)
	last1 := map[int64]int{0: 0}
	last2 := map[int64]int{0: 0}
	lastAll := map[int64]int{0: 0}
	lastPair := map[[2]int64]int{{0, 0}: 0}
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix1 += a[i-1]
		prefix2 += b[i-1]
		prefixAll += a[i-1] + b[i-1]
		best := dp[i-1]
		if j, ok := last1[prefix1]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := last2[prefix2]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := lastAll[prefixAll]; ok && best < dp[j]+1 {
			best = dp[j] + 1
		}
		if j, ok := lastPair[[2]int64{prefix1, prefix2}]; ok && best < dp[j]+2 {
			best = dp[j] + 2
		}
		dp[i] = best
		last1[prefix1] = i
		last2[prefix2] = i
		lastAll[prefixAll] = i
		lastPair[[2]int64{prefix1, prefix2}] = i
	}
	return dp[n]
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			a[i] = v
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			b[i] = v
		}
		expected := solveCase(a, b)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != expected {
			fmt.Printf("case %d failed: expected %d got %s\n", caseIdx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func expected(n, k int, arr []int64) int64 {
	const Inf int64 = math.MaxInt64 / 4
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, k+1)
		for j := range dp[i] {
			dp[i][j] = Inf
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		val := arr[i-1]
		for j := 0; j <= k; j++ {
			if dp[i-1][j]+val < dp[i][j] {
				dp[i][j] = dp[i-1][j] + val
			}
		}
		minVal := val
		for length := 2; length <= k+1 && length <= i; length++ {
			if arr[i-length] < minVal {
				minVal = arr[i-length]
			}
			cost := length - 1
			for j := cost; j <= k; j++ {
				cand := dp[i-length][j-cost] + int64(length)*minVal
				if cand < dp[i][j] {
					dp[i][j] = cand
				}
			}
		}
	}
	ans := dp[n][0]
	for j := 1; j <= k; j++ {
		if dp[n][j] < ans {
			ans = dp[n][j]
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("failed to open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[2+i], 10, 64)
			arr[i] = v
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		want := strconv.FormatInt(expected(n, k, arr), 10)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, want, strings.TrimSpace(got), sb.String())
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

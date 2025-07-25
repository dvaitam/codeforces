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

func expected(a []int) string {
	const inf = int(1e9)
	n := len(a)
	dp := make([][3]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 0; j < 3; j++ {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		dp[i][0] = min3(dp[i-1][0], dp[i-1][1], dp[i-1][2]) + 1
		if a[i-1] == 1 || a[i-1] == 3 {
			dp[i][1] = min(dp[i-1][0], dp[i-1][2])
		}
		if a[i-1] == 2 || a[i-1] == 3 {
			dp[i][2] = min(dp[i-1][0], dp[i-1][1])
		}
	}
	ans := min3(dp[n][0], dp[n][1], dp[n][2])
	return strconv.Itoa(ans)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func min3(a, b, c int) int { return min(a, min(b, c)) }

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || n != len(fields)-1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[i+1])
			nums[i] = v
		}
		want := expected(nums)
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

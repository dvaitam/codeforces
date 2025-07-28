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

func solve(n int, mod int64) int64 {
	dp := make([]int64, n+2)
	diff := make([]int64, n+2)
	dp[1] = 1 % mod
	prefix := dp[1]
	for j := 2; j <= n; j++ {
		l := j
		r := l + j
		if r > n+1 {
			r = n + 1
		}
		diff[l] = (diff[l] + dp[1]) % mod
		diff[r] = (diff[r] - dp[1]) % mod
	}
	cur := int64(0)
	for i := 2; i <= n; i++ {
		cur = (cur + diff[i]) % mod
		if cur < 0 {
			cur += mod
		}
		dp[i] = (prefix + cur) % mod
		prefix = (prefix + dp[i]) % mod
		for j := 2; i*j <= n; j++ {
			l := i * j
			r := l + j
			if r > n+1 {
				r = n + 1
			}
			diff[l] = (diff[l] + dp[i]) % mod
			diff[r] = (diff[r] - dp[i]) % mod
		}
	}
	return dp[n] % mod
}

func runCase(bin string, n int, mod int64) error {
	input := fmt.Sprintf("%d %d\n", n, mod)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprintf("%d", solve(n, mod))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		m64, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		if err := runCase(bin, n, m64); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func expectedD1(n, mod int) int {
	dp := make([]int, n+2)
	pref := make([]int, n+2)
	diff := make([]int, n+2)

	dp[1] = 1
	pref[1] = 1
	for j := 2; j <= n; j++ {
		start := j
		end := j + j
		if end > n+1 {
			end = n + 1
		}
		diff[start] = (diff[start] + dp[1]) % mod
		if end <= n {
			diff[end] = (diff[end] - dp[1]) % mod
		}
	}

	add := 0
	for i := 2; i <= n; i++ {
		add = (add + diff[i]) % mod
		if add < 0 {
			add += mod
		}
		dp[i] = (pref[i-1] + add) % mod
		pref[i] = (pref[i-1] + dp[i]) % mod

		for j := 2; i*j <= n; j++ {
			start := i * j
			end := start + j
			if end > n+1 {
				end = n + 1
			}
			diff[start] = (diff[start] + dp[i]) % mod
			if end <= n {
				diff[end] = (diff[end] - dp[i]) % mod
			}
		}
	}

	ans := dp[n] % mod
	if ans < 0 {
		ans += mod
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
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesD1.txt")
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
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "case %d bad format\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		mod, _ := strconv.Atoi(fields[1])
		expected := expectedD1(n, mod)
		input := fmt.Sprintf("%d %d\n", n, mod)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

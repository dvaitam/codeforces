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

func minBackward(s string) int {
	n := len(s)
	arr := []byte(s)
	full := 1<<n - 1
	dp := make([][]int, 1<<n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = 1 << 30
		}
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = 0
	}
	for mask := 0; mask <= full; mask++ {
		for last := 0; last < n; last++ {
			val := dp[mask][last]
			if val == 1<<30 {
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) != 0 {
					continue
				}
				if arr[next] == arr[last] {
					continue
				}
				nm := mask | (1 << next)
				cost := val
				if next < last {
					cost++
				}
				if cost < dp[nm][next] {
					dp[nm][next] = cost
				}
			}
		}
	}
	best := 1 << 30
	for i := 0; i < n; i++ {
		if dp[full][i] < best {
			best = dp[full][i]
		}
	}
	return best
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		expect := minBackward(s)
		input := fmt.Sprintf("%s\n", s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d non-integer output %s\n", idx, got)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, val)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

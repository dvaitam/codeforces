package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const mod = 1000000007

func solveC(n, m int, s string) int {
	delta := 0
	minPref := 0
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < minPref {
			minPref = bal
		}
	}
	delta = bal
	maxLen := n - m
	dp := make([][]int64, maxLen+1)
	for i := range dp {
		dp[i] = make([]int64, maxLen+1)
	}
	dp[0][0] = 1
	for i := 0; i < maxLen; i++ {
		for j := 0; j <= maxLen; j++ {
			v := dp[i][j]
			if v == 0 {
				continue
			}
			if j+1 <= maxLen {
				dp[i+1][j+1] = (dp[i+1][j+1] + v) % mod
			}
			if j > 0 {
				dp[i+1][j-1] = (dp[i+1][j-1] + v) % mod
			}
		}
	}
	ans := int64(0)
	for l := 0; l <= maxLen; l++ {
		qlen := maxLen - l
		for j := 0; j <= maxLen; j++ {
			v := dp[l][j]
			if v == 0 {
				continue
			}
			if j+minPref < 0 {
				continue
			}
			y := j + delta
			if y < 0 || y > maxLen {
				continue
			}
			add := dp[qlen][y]
			if add == 0 {
				continue
			}
			ans = (ans + v*add) % mod
		}
	}
	return int(ans)
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = '('
		} else {
			b[i] = ')'
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(30) + 1
		m := rand.Intn(n) + 1
		s := randomString(m)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		fmt.Fprintln(&input, s)
		expected := solveC(n, m, s)
		cmd := exec.Command(binary)
		cmd.Stdin = &input
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: binary error: %v\n", t, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", t)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

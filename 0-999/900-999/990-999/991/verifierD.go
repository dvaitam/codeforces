package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expectedD(s1, s2 string) int {
	n := len(s1)
	blocked := make([]int, n)
	for i := 0; i < n; i++ {
		if s1[i] == 'X' {
			blocked[i] |= 1
		}
		if s2[i] == 'X' {
			blocked[i] |= 2
		}
	}
	dp := make([][4]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 0; j < 4; j++ {
			dp[i][j] = -1000000
		}
	}
	dp[0][0] = 0
	for i := 0; i < n; i++ {
		for mask := 0; mask < 4; mask++ {
			if dp[i][mask] < 0 {
				continue
			}
			curBlock := blocked[i]
			dp[i+1][0] = max(dp[i+1][0], dp[i][mask])
			if i+1 >= n {
				continue
			}
			nxtBlock := blocked[i+1]
			if (mask&2) == 0 && (curBlock&2) == 0 && (nxtBlock&1) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][3] = max(dp[i+1][3], dp[i][mask]+1)
			}
			if (mask&1) == 0 && (curBlock&1) == 0 && (mask&2) == 0 && (curBlock&2) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][2] = max(dp[i+1][2], dp[i][mask]+1)
			}
			if (mask&1) == 0 && (curBlock&1) == 0 && (nxtBlock&1) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][3] = max(dp[i+1][3], dp[i][mask]+1)
			}
			if (mask&1) == 0 && (curBlock&1) == 0 && (nxtBlock&1) == 0 && (mask&2) == 0 && (curBlock&2) == 0 {
				dp[i+1][1] = max(dp[i+1][1], dp[i][mask]+1)
			}
		}
	}
	return dp[n][0]
}

func generateCaseD(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	var s1, s2 strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s1.WriteByte('X')
		} else {
			s1.WriteByte('0')
		}
		if rng.Intn(2) == 0 {
			s2.WriteByte('X')
		} else {
			s2.WriteByte('0')
		}
	}
	input := s1.String() + "\n" + s2.String() + "\n"
	expected := expectedD(s1.String(), s2.String())
	return input, expected
}

func runCaseD(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(outStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCaseD(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s, t string) int {
	n, m := len(s), len(t)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + int(s[i-1]-'a'+1)
	}
	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + int(t[j-1]-'a'+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			costDel := dp[i-1][j] + int(s[i-1]-'a'+1)
			costIns := dp[i][j-1] + int(t[j-1]-'a'+1)
			diff := int(s[i-1] - t[j-1])
			if diff < 0 {
				diff = -diff
			}
			costSub := dp[i-1][j-1] + diff
			best := costDel
			if costIns < best {
				best = costIns
			}
			if costSub < best {
				best = costSub
			}
			dp[i][j] = best
		}
	}
	return dp[n][m]
}

func randString(rng *rand.Rand, n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		s := randString(rng, rng.Intn(5)+1)
		t := randString(rng, rng.Intn(5)+1)
		input := fmt.Sprintf("%s\n%s\n", s, t)
		want := expected(s, t)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

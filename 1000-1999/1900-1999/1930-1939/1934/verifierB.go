package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

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
	if errBuf.Len() > 0 {
		return "", fmt.Errorf(errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int) int {
	coins := []int{1, 3, 6, 10, 15}
	const maxBase = 29
	const inf = int(1e9)
	dp := make([]int, maxBase+1)
	for i := 1; i <= maxBase; i++ {
		dp[i] = inf
		for _, c := range coins {
			if i >= c && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	if n < 15 {
		return dp[n]
	}
	r := (n-15)%15 + 15
	ans := (n - r) / 15
	return ans + dp[r]
}

func genTests() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]int, 102)
	for i := 0; i < 102; i++ {
		tests[i] = rng.Intn(1_000_000_000) + 1
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for idx, n := range cases {
		input := fmt.Sprintf("1\n%d\n", n)
		expect := fmt.Sprintf("%d", expected(n))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed: n=%d expected %s got %s\n", idx+1, n, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

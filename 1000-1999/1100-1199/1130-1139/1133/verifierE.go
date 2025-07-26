package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func solveE(n, k int, arr []int) int {
	sort.Ints(arr)
	best := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := i; j >= 1; j-- {
			if arr[i-1]-arr[j-1] > 5 {
				break
			}
			best[i] = j
		}
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	for i := 1; i <= n; i++ {
		copy(dp[i], dp[i-1])
		l := best[i]
		sz := i - l + 1
		if dp[i][1] < sz {
			dp[i][1] = sz
		}
		for j := 2; j <= k; j++ {
			if dp[l-1][j-1]+sz > dp[i][j] {
				dp[i][j] = dp[l-1][j-1] + sz
			}
		}
	}
	ans := 0
	for j := 1; j <= k; j++ {
		if dp[n][j] > ans {
			ans = dp[n][j]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(30)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveE(n, k, arr)
}

func run(bin, input string) (int, error) {
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
		return 0, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	val, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return 0, fmt.Errorf("invalid integer output: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

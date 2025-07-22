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

const mod = 1000000007

func run(bin, input string) (string, error) {
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

func expected(parents []int) int {
	n := len(parents) - 1
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i]
		children[p] = append(children[p], i)
	}
	dp := make([]int64, n+1)
	for i := n; i >= 1; i-- {
		if len(children[i]) == 0 {
			dp[i] = 3
		} else {
			prod := int64(1)
			for _, c := range children[i] {
				prod = prod * dp[c] % mod
			}
			dp[i] = (prod + 2) % mod
		}
	}
	return int(dp[1] % mod)
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(49) + 2 // 2..50
	parents := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parents[i] = rng.Intn(i-1) + 1
	}
	return parents
}

func check(parents []int, out string) error {
	want := expected(parents)
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("invalid integer output")
	}
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
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
		parents := generateCase(rng)
		var sb strings.Builder
		n := len(parents) - 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 2; j <= n; j++ {
			if j > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parents[j]))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(parents, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

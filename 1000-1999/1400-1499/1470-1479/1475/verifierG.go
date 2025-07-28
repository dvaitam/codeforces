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

func solveG(n int, arr []int) string {
	maxA := 0
	for _, v := range arr {
		if v > maxA {
			maxA = v
		}
	}
	freq := make([]int, maxA+1)
	for _, v := range arr {
		freq[v]++
	}
	dp := make([]int, maxA+1)
	for i := 1; i <= maxA; i++ {
		dp[i] = freq[i]
	}
	for i := 1; i <= maxA; i++ {
		if freq[i] == 0 {
			continue
		}
		for j := i * 2; j <= maxA; j += i {
			if freq[j] > 0 && dp[j] < dp[i]+freq[j] {
				dp[j] = dp[i] + freq[j]
			}
		}
	}
	best := 0
	for i := 1; i <= maxA; i++ {
		if dp[i] > best {
			best = dp[i]
		}
	}
	ans := n - best
	return fmt.Sprintf("%d", ans)
}

func generateG(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	out := solveG(n, arr)
	return sb.String(), out
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateG(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

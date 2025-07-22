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

func expected(n, k int, h []int) int {
	sum := 0
	for i := 0; i < k; i++ {
		sum += h[i]
	}
	minSum := sum
	minPos := 0
	for i := k; i < n; i++ {
		sum += h[i] - h[i-k]
		if sum < minSum {
			minSum = sum
			minPos = i - k + 1
		}
	}
	return minPos + 1
}

func generateCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	h := make([]int, n)
	for i := range h {
		h[i] = rng.Intn(100) + 1
	}
	return n, k, h
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, h := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j, v := range h {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		want := fmt.Sprint(expected(n, k, h))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

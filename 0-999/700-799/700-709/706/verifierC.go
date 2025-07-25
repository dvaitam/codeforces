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

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func solveCase(costs []int64, words []string) int64 {
	n := len(costs)
	sr := make([]string, n)
	for i := range words {
		sr[i] = reverse(words[i])
	}
	const INF int64 = 1 << 60
	prev0, prev1 := int64(0), costs[0]
	for i := 1; i < n; i++ {
		cur0, cur1 := INF, INF
		if prev0 < INF && words[i-1] <= words[i] {
			if prev0 < cur0 {
				cur0 = prev0
			}
		}
		if prev1 < INF && sr[i-1] <= words[i] {
			if prev1 < cur0 {
				cur0 = prev1
			}
		}
		if prev0 < INF && words[i-1] <= sr[i] {
			if prev0+costs[i] < cur1 {
				cur1 = prev0 + costs[i]
			}
		}
		if prev1 < INF && sr[i-1] <= sr[i] {
			if prev1+costs[i] < cur1 {
				cur1 = prev1 + costs[i]
			}
		}
		prev0, prev1 = cur0, cur1
	}
	res := prev0
	if prev1 < res {
		res = prev1
	}
	if res >= INF {
		return -1
	}
	return res
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(6) + 1
	costs := make([]int64, n)
	words := make([]string, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := range costs {
		costs[i] = int64(rng.Intn(10) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", costs[i])
	}
	sb.WriteByte('\n')
	for i := range words {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := range b {
			b[j] = byte('a' + rng.Intn(3))
		}
		words[i] = string(b)
		sb.WriteString(words[i])
		if i < n-1 {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	res := solveCase(costs, words)
	return sb.String(), res
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func longestSeq(h []int64, d int64) []int {
	n := len(h)
	dp := make([]int, n)
	prev := make([]int, n)
	best := 0
	bestIdx := 0
	for i := 0; i < n; i++ {
		dp[i] = 1
		prev[i] = -1
		for j := 0; j < i; j++ {
			if abs64(h[i]-h[j]) >= d && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}
		if dp[i] > best {
			best = dp[i]
			bestIdx = i
		}
	}
	seq := make([]int, 0, best)
	for idx := bestIdx; idx != -1; idx = prev[idx] {
		seq = append(seq, idx+1)
	}
	for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
		seq[i], seq[j] = seq[j], seq[i]
	}
	return seq
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func genCase(rng *rand.Rand) (int64, []int64) {
	n := rng.Intn(10) + 1
	d := int64(rng.Intn(5))
	h := make([]int64, n)
	for i := range h {
		h[i] = int64(rng.Intn(50) + 1)
	}
	return d, h
}

func runCase(bin string, d int64, h []int64, exp []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(h), d)
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	var k int
	fmt.Sscan(fields[0], &k)
	if k != len(exp) {
		return fmt.Errorf("expected length %d got %d", len(exp), k)
	}
	if k == 0 {
		return nil
	}
	if len(fields[1:]) != k {
		return fmt.Errorf("expected %d indices got %d", k, len(fields[1:]))
	}
	idxs := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Sscan(fields[1+i], &idxs[i])
	}
	prevIdx := 0
	for i, idx := range idxs {
		if idx < 1 || idx > len(h) {
			return fmt.Errorf("index %d out of range", idx)
		}
		if i > 0 && idx <= prevIdx {
			return fmt.Errorf("indices not increasing")
		}
		if i > 0 {
			if abs64(h[idx-1]-h[idxs[i-1]-1]) < d {
				return fmt.Errorf("invalid jump")
			}
		}
		prevIdx = idx
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		d, h := genCase(rng)
		exp := longestSeq(h, d)
		if err := runCase(bin, d, h, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "d=%d heights=%v\n", d, h)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

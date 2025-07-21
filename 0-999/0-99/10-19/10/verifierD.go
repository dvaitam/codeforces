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

func solveD(a, b []int) (int, []int) {
	n := len(a)
	m := len(b)
	dp := make([]int, m)
	prev := make([]int, m)
	for j := range prev {
		prev[j] = -1
	}
	for i := 0; i < n; i++ {
		current := 0
		last := -1
		for j := 0; j < m; j++ {
			if a[i] == b[j] {
				if current+1 > dp[j] {
					dp[j] = current + 1
					prev[j] = last
				}
			} else if a[i] > b[j] {
				if dp[j] > current {
					current = dp[j]
					last = j
				}
			}
		}
	}
	length := 0
	endIdx := -1
	for j := 0; j < m; j++ {
		if dp[j] > length {
			length = dp[j]
			endIdx = j
		}
	}
	seq := make([]int, 0, length)
	for idx := endIdx; idx != -1; idx = prev[idx] {
		seq = append(seq, b[idx])
	}
	for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
		seq[i], seq[j] = seq[j], seq[i]
	}
	return length, seq
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := make([]int, n)
	b := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(20)
	}
	for i := range b {
		b[i] = rng.Intn(20)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, v := range b {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < m {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	length, seq := solveD(a, b)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d\n", length))
	for i, v := range seq {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	if len(seq) > 0 {
		exp.WriteByte('\n')
	} else {
		exp.WriteString("\n")
	}
	return sb.String(), exp.String()
}

func runCaseD(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.ReplaceAll(out.String(), "\r", "")
	if strings.TrimSpace(outStr) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCaseD(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

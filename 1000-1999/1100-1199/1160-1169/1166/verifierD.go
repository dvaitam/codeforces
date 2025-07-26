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

func tryMake(a, b, m int64, k int, p []int64) ([]int64, bool) {
	b2 := b - p[k]*a
	if b2 < p[k] || (b2+p[k]-1)/p[k] > m {
		return nil, false
	}
	r := make([]int64, k+1)
	for i := 1; i < k; i++ {
		r[i] = 1
		b2 -= p[i]
	}
	for i := 1; i < k; i++ {
		if b2 <= 0 {
			break
		}
		pi := p[k-i]
		maxAdd := b2 / pi
		if maxAdd > m-1 {
			maxAdd = m - 1
		}
		r[i] += maxAdd
		b2 -= pi * maxAdd
	}
	if b2 != 0 {
		return nil, false
	}
	return r, true
}

func expectedSingle(a, b, m int64) string {
	if a == b {
		return fmt.Sprintf("1 %d", a)
	}
	const K = 51
	p := make([]int64, K)
	p[1], p[2] = 1, 1
	for i := 3; i < K; i++ {
		p[i] = p[i-1] * 2
	}
	for k := 2; k < K; k++ {
		if p[k]*a > b {
			break
		}
		r, ok := tryMake(a, b, m, k, p)
		if !ok {
			continue
		}
		seq := make([]int64, k)
		x := a
		var sum int64
		for i := 0; i < k; i++ {
			seq[i] = x
			sum += x
			x = sum + r[i+1]
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d ", k))
		for i, v := range seq {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		return sb.String()
	}
	return "-1"
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	var expLines []string
	for i := 0; i < t; i++ {
		a := int64(rng.Intn(20) + 1)
		b := a + int64(rng.Intn(50)+1)
		m := int64(rng.Intn(10) + 1)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, m))
		expLines = append(expLines, expectedSingle(a, b, m))
	}
	return sb.String(), strings.Join(expLines, "\n")
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

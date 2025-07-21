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

type interval struct{ l, r int }

func solve(n, k int, h []int) (int, []interval) {
	maxDeque := make([]int, 0, n)
	minDeque := make([]int, 0, n)
	l := 0
	best := 0
	res := make([]interval, 0)
	for r := 0; r < n; r++ {
		for len(maxDeque) > 0 && h[maxDeque[len(maxDeque)-1]] <= h[r] {
			maxDeque = maxDeque[:len(maxDeque)-1]
		}
		maxDeque = append(maxDeque, r)
		for len(minDeque) > 0 && h[minDeque[len(minDeque)-1]] >= h[r] {
			minDeque = minDeque[:len(minDeque)-1]
		}
		minDeque = append(minDeque, r)
		for l <= r && h[maxDeque[0]]-h[minDeque[0]] > k {
			if maxDeque[0] == l {
				maxDeque = maxDeque[1:]
			}
			if minDeque[0] == l {
				minDeque = minDeque[1:]
			}
			l++
		}
		length := r - l + 1
		if length > best {
			best = length
			res = res[:0]
			res = append(res, interval{l, r})
		} else if length == best {
			res = append(res, interval{l, r})
		}
	}
	return best, res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	k := rng.Intn(1000)
	h := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		h[i] = rng.Intn(1000) + 1
		sb.WriteString(fmt.Sprintf("%d", h[i]))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	best, res := solve(n, k, h)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d %d\n", best, len(res)))
	for _, iv := range res {
		exp.WriteString(fmt.Sprintf("%d %d\n", iv.l+1, iv.r+1))
	}
	return sb.String(), strings.TrimSpace(exp.String())
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

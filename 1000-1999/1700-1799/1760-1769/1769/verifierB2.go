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

func ceilDiv(a, b int64) int64 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

func expected(a []int64) []int {
	var total int64
	for _, v := range a {
		total += v
	}
	found := make([]bool, 101)
	var prefix int64
	for _, ai := range a {
		for p := 0; p <= 100; p++ {
			l1 := ceilDiv(int64(p)*ai, 100)
			r1 := (int64(p+1)*ai - 1) / 100
			if l1 > r1 {
				continue
			}
			l2 := ceilDiv(int64(p)*total, 100) - prefix
			r2 := (int64(p+1)*total-1)/100 - prefix
			if l2 > r2 {
				continue
			}
			l := l1
			if l2 > l {
				l = l2
			}
			if l < 0 {
				l = 0
			}
			r := r1
			if r2 < r {
				r = r2
			}
			if r > ai {
				r = ai
			}
			if l <= r {
				found[p] = true
			}
		}
		prefix += ai
	}
	res := make([]int, 0)
	for p := 0; p <= 100; p++ {
		if found[p] {
			res = append(res, p)
		}
	}
	return res
}

func runCase(bin string, a []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	want := expected(a)
	if len(fields) != len(want) {
		return fmt.Errorf("expected %d numbers got %d", len(want), len(fields))
	}
	for i, f := range fields {
		var x int
		if _, err := fmt.Sscan(f, &x); err != nil {
			return fmt.Errorf("invalid int output: %v", err)
		}
		if x != want[i] {
			return fmt.Errorf("mismatch at pos %d: expected %d got %d", i, want[i], x)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) []int64 {
	n := rng.Intn(20) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1e10) + 1
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

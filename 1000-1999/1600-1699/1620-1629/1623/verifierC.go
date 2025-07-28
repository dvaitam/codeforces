package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func canAchieve(h []int64, target int64) bool {
	n := len(h)
	b := make([]int64, n)
	copy(b, h)
	for i := n - 1; i >= 2; i-- {
		if b[i] < target {
			return false
		}
		d := (b[i] - target) / 3
		if d > h[i]/3 {
			d = h[i] / 3
		}
		b[i-1] += d
		b[i-2] += 2 * d
	}
	return b[0] >= target && b[1] >= target
}

func expected(h []int64) int64 {
	l, r := int64(0), int64(1e9)
	for l < r {
		mid := (l + r + 1) / 2
		if canAchieve(h, mid) {
			l = mid
		} else {
			r = mid - 1
		}
	}
	return l
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 3
		h := make([]int64, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			h[i] = rng.Int63n(1000) + 1
			sb.WriteString(fmt.Sprintf("%d ", h[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := fmt.Sprintf("%d", expected(h))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

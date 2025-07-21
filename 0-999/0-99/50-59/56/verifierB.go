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

func expectedAnswer(p []int) (int, int) {
	n := len(p) - 1
	l := 1
	for l <= n && p[l] == l {
		l++
	}
	if l > n {
		return 0, 0
	}
	r := n
	for r >= l && p[r] != l {
		r--
	}
	if r <= l {
		return 0, 0
	}
	for i := l; i <= r; i++ {
		if p[i] != l+r-i {
			return 0, 0
		}
	}
	for i := 1; i < l; i++ {
		if p[i] != i {
			return 0, 0
		}
	}
	for i := r + 1; i <= n; i++ {
		if p[i] != i {
			return 0, 0
		}
	}
	return l, r
}

func generateCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(20) + 2 // at least 2
	perm := make([]int, n+1)
	for i := 1; i <= n; i++ {
		perm[i] = i
	}
	if rng.Intn(2) == 0 {
		l := rng.Intn(n-1) + 1
		r := rng.Intn(n-l) + l + 1
		for i, j := l, r; i < j; i, j = i+1, j-1 {
			perm[i], perm[j] = perm[j], perm[i]
		}
	} else {
		// random permutation not necessarily single reversal
		for i := n; i >= 1; i-- {
			j := rng.Intn(i) + 1
			perm[i], perm[j] = perm[j], perm[i]
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "%d ", perm[i])
	}
	b.WriteByte('\n')
	l, r := expectedAnswer(perm)
	return b.String(), l, r
}

func runCase(bin, input string, l, r int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != 2 {
		return fmt.Errorf("expected two numbers, got: %s", out.String())
	}
	gotL, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("bad output: %s", parts[0])
	}
	gotR, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("bad output: %s", parts[1])
	}
	if gotL != l || gotR != r {
		return fmt.Errorf("expected %d %d got %d %d", l, r, gotL, gotR)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, l, r := generateCase(rng)
		if err := runCase(bin, input, l, r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

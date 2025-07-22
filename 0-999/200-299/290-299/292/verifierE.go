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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1 // 1..8
	m := rng.Intn(20) + 1
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(rng.Intn(41) - 20)
	}
	for i := 1; i <= n; i++ {
		b[i] = int64(rng.Intn(41) - 20)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d", a[i])
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d", b[i])
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	var out strings.Builder
	for i := 0; i < m; i++ {
		t := rng.Intn(2) + 1
		if t == 1 {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			maxk := n - max(x, y) + 1
			if maxk <= 0 {
				maxk = 1
			}
			k := rng.Intn(maxk) + 1
			fmt.Fprintf(&sb, "1 %d %d %d\n", x, y, k)
			for j := 0; j < k; j++ {
				b[y+j] = a[x+j]
			}
		} else {
			x := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "2 %d\n", x)
			fmt.Fprintf(&out, "%d\n", b[x])
		}
	}
	return sb.String(), strings.TrimSpace(out.String())
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
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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

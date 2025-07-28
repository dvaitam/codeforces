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

func run(bin, input string) (string, error) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(50) + 2
	k := rng.Intn(20) + 2
	h := make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = rng.Intn(50)
	}
	return n, k, h
}

func expected(n, k int, h []int) string {
	low, high := h[0], h[0]
	for i := 1; i < n; i++ {
		l := max(h[i], low-(k-1))
		r := min(h[i]+k-1, high+(k-1))
		if l > r {
			return "NO"
		}
		low, high = l, r
	}
	if h[n-1] >= low && h[n-1] <= high {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, n, k int, h []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got := strings.ToUpper(strings.TrimSpace(out))
	expect := expected(n, k, h)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, h := generateCase(rng)
		if err := runCase(bin, n, k, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

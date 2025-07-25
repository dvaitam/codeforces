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

func expected(n, m, k int, a []int, b []int) string {
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pos[a[i-1]] = i
	}
	result := int64(0)
	for _, app := range b {
		p := pos[app]
		screen := (p-1)/k + 1
		result += int64(screen)
		if p > 1 {
			prev := a[p-2]
			a[p-2], a[p-1] = a[p-1], a[p-2]
			pos[app] = p - 1
			pos[prev] = p
		}
	}
	return fmt.Sprintf("%d", result)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	if m > 20 {
		m = 20
	}
	k := rng.Intn(n) + 1
	// permutation a
	perm := rng.Perm(n)
	a := make([]int, n)
	for i, v := range perm {
		a[i] = v + 1
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	exp := expected(n, m, k, append([]int(nil), a...), append([]int(nil), b...))
	return input, exp
}

func runCase(bin, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
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

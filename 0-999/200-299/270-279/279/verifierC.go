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

func solveC(n, m int, a []int, queries [][2]int) string {
	p := make([]int, n)
	p[n-1] = n - 1
	for i := n - 2; i >= 0; i-- {
		if a[i] <= a[i+1] {
			p[i] = p[i+1]
		} else {
			p[i] = i
		}
	}
	q := make([]int, n)
	q[n-1] = n - 1
	for i := n - 2; i >= 0; i-- {
		if a[i] >= a[i+1] {
			q[i] = q[i+1]
		} else {
			q[i] = i
		}
	}
	var sb strings.Builder
	for idx, qr := range queries {
		l, r := qr[0]-1, qr[1]-1
		x := p[l]
		if x >= r || q[x] >= r {
			sb.WriteString("Yes")
		} else {
			sb.WriteString("No")
		}
		if idx+1 < len(queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	m := rng.Intn(15) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(20)
	}
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, qr := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qr[0], qr[1])
	}
	input := sb.String()
	expected := solveC(n, m, a, queries)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		in, exp := generateCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

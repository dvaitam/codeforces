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

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func expected(n, m, a0, b0 int) string {
	switch {
	case n > m:
		if (a0 > 0 && b0 > 0) || (a0 < 0 && b0 < 0) {
			return "Infinity"
		}
		return "-Infinity"
	case n < m:
		return "0/1"
	default:
		num := a0
		den := b0
		if den < 0 {
			num = -num
			den = -den
		}
		g := gcd(num, den)
		num /= g
		den /= g
		return fmt.Sprintf("%d/%d", num, den)
	}
}

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

func buildInput(rng *rand.Rand, n, m, a0, b0 int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	sb.WriteString(fmt.Sprintf("%d", a0))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf(" %d", rng.Intn(201)-100))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d", b0))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf(" %d", rng.Intn(201)-100))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type caseB struct {
		n, m   int
		a0, b0 int
	}
	var cases []caseB
	edge := []caseB{
		{2, 1, 1, 1},
		{3, 2, -2, 5},
		{1, 2, 10, -5},
		{2, 2, 4, 2},
		{1, 1, -2, 4},
		{1, 1, 2, -3},
		{0, 0, -1, -1},
		{0, 0, 5, -10},
		{5, 2, 100, 100},
		{5, 5, 5, 7},
	}
	cases = append(cases, edge...)
	for i := 0; i < 100; i++ {
		n := rng.Intn(101)
		m := rng.Intn(101)
		a0 := rng.Intn(201) - 100
		if a0 == 0 {
			a0 = 1
		}
		b0 := rng.Intn(201) - 100
		if b0 == 0 {
			b0 = 1
		}
		cases = append(cases, caseB{n, m, a0, b0})
	}

	for i, c := range cases {
		input := buildInput(rng, c.n, c.m, c.a0, c.b0)
		want := expected(c.n, c.m, c.a0, c.b0)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

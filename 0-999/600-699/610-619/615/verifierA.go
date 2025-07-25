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

func expected(n, m int, sets [][]int) string {
	bulbs := make([]bool, m+1)
	for _, s := range sets {
		for _, b := range s {
			if b >= 1 && b <= m {
				bulbs[b] = true
			}
		}
	}
	for i := 1; i <= m; i++ {
		if !bulbs[i] {
			return "NO"
		}
	}
	return "YES"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	sets := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		x := rng.Intn(m + 1)
		sb.WriteString(fmt.Sprintf("%d", x))
		sets[i] = make([]int, x)
		for j := 0; j < x; j++ {
			y := rng.Intn(m) + 1
			sets[i][j] = y
			sb.WriteString(fmt.Sprintf(" %d", y))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, m, sets)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []string{
		"3 4\n1 1\n1 2\n2 3 4\n",
		"2 2\n1 1\n0\n",
	}
	exps := []string{"YES", "NO"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}

	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

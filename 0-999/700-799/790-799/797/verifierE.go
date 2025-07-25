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

func expectedE(a []int, p, k int) int {
	cnt := 0
	n := len(a) - 1
	for p <= n {
		p = p + a[p] + k
		cnt++
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, []int, [][2]int) {
	n := rng.Intn(20) + 1
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(n) + 1
	}
	q := rng.Intn(20) + 1
	queries := make([][2]int, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		p := rng.Intn(n) + 1
		k := rng.Intn(n) + 1
		queries[i] = [2]int{p, k}
		sb.WriteString(fmt.Sprintf("%d %d\n", p, k))
	}
	return sb.String(), a, queries
}

func runCase(bin string, input string, a []int, queries [][2]int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(resFields) != len(queries) {
		return fmt.Errorf("expected %d outputs got %d (%q)", len(queries), len(resFields), out.String())
	}
	for i, f := range resFields {
		got, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad output %q", f)
		}
		exp := expectedE(a, queries[i][0], queries[i][1])
		if got != exp {
			return fmt.Errorf("query %d expected %d got %d", i+1, exp, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, a, qs := generateCase(rng)
		if err := runCase(bin, in, a, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

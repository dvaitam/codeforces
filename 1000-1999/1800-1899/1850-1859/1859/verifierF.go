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

func runReference(input string) (string, error) {
	cmd := exec.Command("go", "run", "1859F.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type edge struct{ u, v, w int }

func generateCase(rng *rand.Rand) (int, int, []edge, string, [][2]int) {
	n := rng.Intn(5) + 2
	T := rng.Intn(5) + 1
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		w := rng.Intn(10) + 1
		edges[i-2] = edge{u, i, w}
	}
	sbytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sbytes[i] = '0'
		} else {
			sbytes[i] = '1'
		}
	}
	s := string(sbytes)
	q := rng.Intn(3) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		queries[i] = [2]int{a, b}
	}
	return n, T, edges, s, queries
}
func buildInput(n, T int, edges []edge, s string, qs [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, T))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(s + "\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for i, q := range qs {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%d %d", q[0], q[1]))
	}
	sb.WriteString("\n")
	return sb.String()
}

func runCase(bin string, n, T int, edges []edge, s string, qs [][2]int) error {
	input := buildInput(n, T, edges, s, qs)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	exp, err := runReference(input)
	if err != nil {
		return fmt.Errorf("reference failed: %v\n%s", err, exp)
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, T, edges, s, qs := generateCase(rng)
		if err := runCase(bin, n, T, edges, s, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

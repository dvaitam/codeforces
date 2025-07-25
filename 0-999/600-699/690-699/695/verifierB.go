package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "695B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runOracle(oracle, input string) (string, error) {
	cmd := exec.Command(oracle)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxEdges-(n-1)+1)
	s := rng.Intn(n) + 1
	t := rng.Intn(n-1) + 1
	if t >= s {
		t++
	}
	type pair struct{ u, v int }
	edges := make([][3]int, 0, m)
	used := make(map[pair]bool)
	for i := 2; i <= n; i++ {
		u := i
		v := rng.Intn(i-1) + 1
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{u, v, w})
		if u > v {
			u, v = v, u
		}
		used[pair{u, v}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		du, dv := u, v
		if du > dv {
			du, dv = dv, du
		}
		if used[pair{du, dv}] {
			continue
		}
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{u, v, w})
		used[pair{du, dv}] = true
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	expected, err := runOracle(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

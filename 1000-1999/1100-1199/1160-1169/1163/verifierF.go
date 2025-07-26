package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1163F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(48))
	var tests []string
	// simple path graph
	tests = append(tests, "2 1 1\n1 2 5\n1 3\n")
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		m := rng.Intn(6) + n - 1
		q := rng.Intn(3) + 1
		edges := make([][3]int, m)
		for i := 0; i < n-1; i++ {
			edges[i] = [3]int{i + 1, i + 2, rng.Intn(10) + 1}
		}
		for i := n - 1; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				v = (v % n) + 1
			}
			w := rng.Intn(10) + 1
			edges[i] = [3]int{u, v, w}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		for i := 0; i < q; i++ {
			ti := rng.Intn(m) + 1
			x := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", ti, x))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	tests := generateTests()
	for i, input := range tests {
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

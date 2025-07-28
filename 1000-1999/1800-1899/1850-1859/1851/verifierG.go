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
	oracle := filepath.Join(os.TempDir(), "oracleG1851.bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1851G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generate() string {
	const T = 100
	rng := rand.New(rand.NewSource(7))
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(n*(n-1)/2) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
		}
		sb.WriteByte('\n')
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			a, b := u, v
			if a > b {
				a, b = b, a
			}
			key := [2]int{a, b}
			if used[key] {
				continue
			}
			used[key] = true
			edges = append(edges, key)
		}
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		q := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for j := 0; j < q; j++ {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			e := rng.Intn(20)
			fmt.Fprintf(&sb, "%d %d %d\n", a, b, e)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	input := generate()
	exp, err := run(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	got, err := run(cand, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		fmt.Println("wrong answer")
		fmt.Println("input:\n" + input)
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

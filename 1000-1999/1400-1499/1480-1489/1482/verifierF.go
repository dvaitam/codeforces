package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildOracle() (string, error) {
	oracle := "oracleF"
	cmd := exec.Command("go", "build", "-o", oracle, "1482F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run error: %v\nstderr: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rand.Seed(6)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 2
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges) + 1
		edges := make([][3]int, 0, m)
		used := map[[2]int]bool{}
		for len(edges) < m {
			u := rand.Intn(n)
			v := rand.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if used[key] {
				continue
			}
			used[key] = true
			w := rand.Intn(5) + 1
			edges = append(edges, [3]int{u + 1, v + 1, w})
		}
		q := rand.Intn(3) + 1
		queries := make([][3]int, q)
		for j := 0; j < q; j++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			l := rand.Intn(5) + 1
			queries[j] = [3]int{u, v, l}
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d %d\n", e[0], e[1], e[2])
		}
		fmt.Fprintln(&input, q)
		for _, qu := range queries {
			fmt.Fprintf(&input, "%d %d %d\n", qu[0], qu[1], qu[2])
		}
		inp := input.String()
		expected, err := run("./"+oracle, inp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d mismatch\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

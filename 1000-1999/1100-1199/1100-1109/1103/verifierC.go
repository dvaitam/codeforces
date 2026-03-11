package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	refSource := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSource == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	oracle := filepath.Join(os.TempDir(), "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

type graphInfo struct {
	n, m, k int
	adj      map[int]map[int]bool // adjacency as set for quick lookup
}

func parseInput(input string) graphInfo {
	tokens := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(tokens[idx])
		idx++
		return v
	}
	n := nextInt()
	m := nextInt()
	k := nextInt()
	adj := make(map[int]map[int]bool)
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]bool)
	}
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		adj[u][v] = true
		adj[v][u] = true
	}
	return graphInfo{n, m, k, adj}
}

func validateOutput(g graphInfo, output string) error {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	idx := 0
	nextStr := func() (string, error) {
		if idx >= len(tokens) {
			return "", fmt.Errorf("unexpected end of output")
		}
		s := tokens[idx]
		idx++
		return s, nil
	}
	nextInt := func() (int, error) {
		s, err := nextStr()
		if err != nil {
			return 0, err
		}
		return strconv.Atoi(s)
	}

	keyword, err := nextStr()
	if err != nil {
		return err
	}

	if keyword == "-1" {
		return fmt.Errorf("candidate returned -1 but a solution should exist")
	}

	switch keyword {
	case "PATH":
		c, err := nextInt()
		if err != nil {
			return fmt.Errorf("cannot read path length: %v", err)
		}
		needed := g.n / g.k
		if g.n%g.k != 0 {
			needed++ // ceil(n/k) -- actually the problem says "at least n/k"
		}
		// The problem says length >= n/k (not necessarily ceiling).
		// "simple path with length at least n/k" where length = number of vertices.
		// Actually re-reading: "length at least n/k" means c >= ceil(n/k)? No:
		// "a simple path with length at least n/k (n is not necessarily divided by k)"
		// This means c >= n/k, and since c is integer, c >= ceil(n/k).
		if c < (g.n+g.k-1)/g.k {
			return fmt.Errorf("path length %d < ceil(%d/%d) = %d", c, g.n, g.k, (g.n+g.k-1)/g.k)
		}
		path := make([]int, c)
		seen := make(map[int]bool)
		for i := 0; i < c; i++ {
			v, err := nextInt()
			if err != nil {
				return fmt.Errorf("cannot read path vertex %d: %v", i, err)
			}
			if v < 1 || v > g.n {
				return fmt.Errorf("path vertex %d out of range", v)
			}
			if seen[v] {
				return fmt.Errorf("path vertex %d repeated", v)
			}
			seen[v] = true
			path[i] = v
		}
		for i := 0; i+1 < c; i++ {
			if !g.adj[path[i]][path[i+1]] {
				return fmt.Errorf("no edge between path vertices %d and %d", path[i], path[i+1])
			}
		}

	case "CYCLES":
		usedRep := make(map[int]bool)
		for ci := 0; ci < g.k; ci++ {
			c, err := nextInt()
			if err != nil {
				return fmt.Errorf("cannot read cycle %d length: %v", ci+1, err)
			}
			if c < 3 {
				return fmt.Errorf("cycle %d length %d < 3", ci+1, c)
			}
			if c%3 == 0 {
				return fmt.Errorf("cycle %d length %d divisible by 3", ci+1, c)
			}
			verts := make([]int, c)
			seen := make(map[int]bool)
			for j := 0; j < c; j++ {
				v, err := nextInt()
				if err != nil {
					return fmt.Errorf("cannot read cycle %d vertex %d: %v", ci+1, j, err)
				}
				if v < 1 || v > g.n {
					return fmt.Errorf("cycle %d vertex %d out of range", ci+1, v)
				}
				if seen[v] {
					return fmt.Errorf("cycle %d vertex %d repeated", ci+1, v)
				}
				seen[v] = true
				verts[j] = v
			}
			// Check edges form a cycle
			for j := 0; j < c; j++ {
				u := verts[j]
				v := verts[(j+1)%c]
				if !g.adj[u][v] {
					return fmt.Errorf("cycle %d: no edge between %d and %d", ci+1, u, v)
				}
			}
			// Representative is the first vertex
			rep := verts[0]
			if usedRep[rep] {
				return fmt.Errorf("cycle %d: representative %d already used", ci+1, rep)
			}
			usedRep[rep] = true
		}

	default:
		return fmt.Errorf("unexpected keyword %q (expected PATH, CYCLES, or -1)", keyword)
	}
	return nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 4 // 4..8
	k := r.Intn(n) + 1
	m := n * (n - 1) / 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", i, j))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	var cases []string
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		// Run oracle to confirm a solution exists (not -1)
		oracleOut, err := run(oracle, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}

		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}

		// If oracle says -1, candidate must also say -1
		if strings.TrimSpace(oracleOut) == "-1" {
			if strings.TrimSpace(got) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d: oracle says -1 but candidate produced output\n", i+1)
				os.Exit(1)
			}
			continue
		}

		g := parseInput(tc)
		if err := validateOutput(g, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\ngot:\n%s\n", i+1, err, tc, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

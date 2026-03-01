package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input string
}

func maxEdges(n, k int) int {
	return (n - k) + ((n-1)*(n-2))/2
}

func impossible(n, m, k int) bool {
	if m < n-1 {
		return true
	}
	if k == n {
		return true
	}
	return m > maxEdges(n, k)
}

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(s)
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func validateOutput(n, m, k int, special []int, out string) error {
	vals, err := parseInts(out)
	if err != nil {
		return fmt.Errorf("output contains non-integer token: %w", err)
	}
	if len(vals) == 1 && vals[0] == -1 {
		if impossible(n, m, k) {
			return nil
		}
		return fmt.Errorf("printed -1 although a construction exists")
	}
	if impossible(n, m, k) {
		return fmt.Errorf("construction is impossible, expected -1")
	}
	if len(vals)%2 != 0 {
		return fmt.Errorf("edge list must contain pairs, got %d integers", len(vals))
	}
	edgesCnt := len(vals) / 2
	if edgesCnt != m {
		return fmt.Errorf("expected %d edges, got %d", m, edgesCnt)
	}

	forbiddenRoot := special[0]
	otherSpecial := make(map[int]bool, k-1)
	for i := 1; i < k; i++ {
		otherSpecial[special[i]] = true
	}

	adj := make([][]int, n+1)
	seen := make(map[[2]int]bool, m)
	for i := 0; i < len(vals); i += 2 {
		u, v := vals[i], vals[i+1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge (%d,%d) out of range 1..%d", u, v, n)
		}
		if u == v {
			return fmt.Errorf("self-loop at %d", u)
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if seen[key] {
			return fmt.Errorf("duplicate edge (%d,%d)", a, b)
		}
		seen[key] = true

		if (u == forbiddenRoot && otherSpecial[v]) || (v == forbiddenRoot && otherSpecial[u]) {
			return fmt.Errorf("forbidden edge between %d and special node", forbiddenRoot)
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// connectedness check
	q := []int{1}
	vis := make([]bool, n+1)
	vis[1] = true
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return fmt.Errorf("graph is not connected")
		}
	}

	return nil
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(8) + 3
		k := rng.Intn(n-1) + 2
		special := make([]int, 0, k)
		used := make(map[int]bool)
		for len(special) < k {
			x := rng.Intn(n) + 1
			if !used[x] {
				used[x] = true
				special = append(special, x)
			}
		}
		m := rng.Intn(maxEdges(n, k) + 2) // includes some impossible cases
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i, v := range special {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{sb.String()})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		parts, err := parseInts(t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error: %v\n", err)
			os.Exit(1)
		}
		n, m, k := parts[0], parts[1], parts[2]
		special := parts[3:]

		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validateOutput(n, m, k, special, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, t.input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

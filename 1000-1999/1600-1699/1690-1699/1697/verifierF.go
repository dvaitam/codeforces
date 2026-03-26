package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---------- embedded reference solver for 1697F (2-SAT) ----------

func referenceSolve(input string) string {
	b := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(b) && b[pos] <= ' ' {
			pos++
		}
		if pos >= len(b) {
			return 0
		}
		neg := false
		if b[pos] == '-' {
			neg = true
			pos++
		}
		res := 0
		for pos < len(b) && b[pos] >= '0' && b[pos] <= '9' {
			res = res*10 + int(b[pos]-'0')
			pos++
		}
		if neg {
			return -res
		}
		return res
	}

	t := nextInt()
	if t == 0 {
		return ""
	}

	var out bytes.Buffer
	w := bufio.NewWriter(&out)

	maxN := 20000*9 + 5
	adj := make([][]int, 2*maxN)
	dfn := make([]int, 2*maxN)
	low := make([]int, 2*maxN)
	scc := make([]int, 2*maxN)
	inStk := make([]bool, 2*maxN)
	stk := make([]int, 0, 2*maxN)

	for tc := 0; tc < t; tc++ {
		n := nextInt()
		m := nextInt()
		k := nextInt()

		nVar := n*(k-1) + 1
		nVertices := 2 * nVar

		for i := 0; i < nVertices; i++ {
			adj[i] = adj[i][:0]
			dfn[i] = 0
			low[i] = 0
			scc[i] = 0
			inStk[i] = false
		}

		addClause := func(u, v int) {
			adj[u^1] = append(adj[u^1], v)
			adj[v^1] = append(adj[v^1], u)
		}

		lit := func(i, v int) int {
			if v <= 1 {
				return 0
			}
			if v >= k+1 {
				return 1
			}
			return 2 * ((i-1)*(k-1) + v - 1)
		}

		addClause(0, 0)

		for i := 1; i <= n; i++ {
			for v := 2; v <= k; v++ {
				addClause(lit(i, v)^1, lit(i, v-1))
			}
		}

		for i := 1; i < n; i++ {
			for v := 2; v <= k; v++ {
				addClause(lit(i, v)^1, lit(i+1, v))
			}
		}

		for idx := 0; idx < m; idx++ {
			typ := nextInt()
			if typ == 1 {
				i := nextInt()
				x := nextInt()
				addClause(lit(i, x)^1, lit(i, x+1))
			} else if typ == 2 {
				i := nextInt()
				j := nextInt()
				x := nextInt()
				for u := 1; u <= k; u++ {
					addClause(lit(i, u)^1, lit(j, x-u+1)^1)
				}
			} else if typ == 3 {
				i := nextInt()
				j := nextInt()
				x := nextInt()
				for u := 1; u <= k; u++ {
					addClause(lit(j, u+1), lit(i, x-u))
				}
			}
		}

		timer := 0
		sccCnt := 0
		stk = stk[:0]

		var tarjan func(u int)
		tarjan = func(u int) {
			timer++
			dfn[u] = timer
			low[u] = timer
			stk = append(stk, u)
			inStk[u] = true
			for _, v := range adj[u] {
				if dfn[v] == 0 {
					tarjan(v)
					if low[v] < low[u] {
						low[u] = low[v]
					}
				} else if inStk[v] {
					if dfn[v] < low[u] {
						low[u] = dfn[v]
					}
				}
			}
			if low[u] == dfn[u] {
				for {
					v := stk[len(stk)-1]
					stk = stk[:len(stk)-1]
					inStk[v] = false
					scc[v] = sccCnt
					if u == v {
						break
					}
				}
				sccCnt++
			}
		}

		for i := 0; i < nVertices; i++ {
			if dfn[i] == 0 {
				tarjan(i)
			}
		}

		possible := true
		for i := 0; i < nVar; i++ {
			if scc[2*i] == scc[2*i+1] {
				possible = false
				break
			}
		}

		if !possible {
			fmt.Fprintln(w, "-1")
		} else {
			for i := 1; i <= n; i++ {
				val := 1
				for v := 2; v <= k; v++ {
					varIdx := (i-1)*(k-1) + v - 1
					if scc[2*varIdx] < scc[2*varIdx+1] {
						val = v
					}
				}
				if i > 1 {
					w.WriteByte(' ')
				}
				fmt.Fprint(w, val)
			}
			w.WriteByte('\n')
		}
	}
	w.Flush()
	return strings.TrimSpace(out.String())
}

// ---------- test generator ----------

type testCase struct {
	input string
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	K := rng.Intn(3) + 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, K))
	for i := 0; i < m; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 {
			idx := rng.Intn(n) + 1
			v := rng.Intn(K) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", idx, v))
		} else if typ == 2 {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			v := rng.Intn(2*K) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", a, b, v))
		} else {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			v := rng.Intn(2*K) + 1
			sb.WriteString(fmt.Sprintf("3 %d %d %d\n", a, b, v))
		}
	}
	return testCase{input: sb.String()}
}

func runBinary(bin, input string) (string, error) {
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

// validate checks that candidate output satisfies all constraints.
// Returns nil if valid, error otherwise.
func validate(input, output string) error {
	b := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(b) && b[pos] <= ' ' {
			pos++
		}
		res := 0
		neg := false
		if pos < len(b) && b[pos] == '-' {
			neg = true
			pos++
		}
		for pos < len(b) && b[pos] >= '0' && b[pos] <= '9' {
			res = res*10 + int(b[pos]-'0')
			pos++
		}
		if neg {
			return -res
		}
		return res
	}

	t := nextInt()
	outLines := strings.Split(strings.TrimSpace(output), "\n")
	li := 0
	for tc := 0; tc < t; tc++ {
		n := nextInt()
		m := nextInt()
		k := nextInt()

		type constraint struct {
			typ        int
			i, j, v   int
		}
		constraints := make([]constraint, m)
		for idx := 0; idx < m; idx++ {
			c := constraint{typ: nextInt()}
			if c.typ == 1 {
				c.i = nextInt()
				c.v = nextInt()
			} else {
				c.i = nextInt()
				c.j = nextInt()
				c.v = nextInt()
			}
			constraints[idx] = c
		}

		if li >= len(outLines) {
			return fmt.Errorf("test %d: missing output line", tc+1)
		}
		line := strings.TrimSpace(outLines[li])
		li++

		if line == "-1" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != n {
			return fmt.Errorf("test %d: expected %d values, got %d", tc+1, n, len(parts))
		}
		vals := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Sscanf(parts[i-1], "%d", &vals[i])
			if vals[i] < 1 || vals[i] > k {
				return fmt.Errorf("test %d: value %d out of range [1,%d]", tc+1, vals[i], k)
			}
		}
		// non-decreasing
		for i := 1; i < n; i++ {
			if vals[i] > vals[i+1] {
				return fmt.Errorf("test %d: not non-decreasing at pos %d", tc+1, i)
			}
		}
		// check constraints
		for _, c := range constraints {
			switch c.typ {
			case 1:
				if vals[c.i] != c.v {
					return fmt.Errorf("test %d: constraint 1 %d %d failed (got %d)", tc+1, c.i, c.v, vals[c.i])
				}
			case 2:
				if vals[c.i]+vals[c.j] > c.v {
					return fmt.Errorf("test %d: constraint 2 %d %d %d failed (%d+%d=%d)", tc+1, c.i, c.j, c.v, vals[c.i], vals[c.j], vals[c.i]+vals[c.j])
				}
			case 3:
				if vals[c.i]+vals[c.j] < c.v {
					return fmt.Errorf("test %d: constraint 3 %d %d %d failed (%d+%d=%d)", tc+1, c.i, c.j, c.v, vals[c.i], vals[c.j], vals[c.i]+vals[c.j])
				}
			}
		}
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

	// Increase stack for recursive tarjan in reference solver
	// by using goroutine with large stack - not needed for small tests

	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		cases[i] = genCase(rng)
	}

	for i, tc := range cases {
		expected := referenceSolve(tc.input)
		got, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		if got == expected {
			continue
		}

		// If outputs differ, both could be valid (multiple solutions).
		// Check if -1 vs non-negative-1 mismatch.
		if (got == "-1") != (expected == "-1") {
			fmt.Fprintf(os.Stderr, "case %d: feasibility mismatch: expected %s got %s\ninput:\n%s", i+1, expected, got, tc.input)
			os.Exit(1)
		}

		// Both non-"-1": validate candidate output against constraints.
		if err := validate(tc.input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
	}
	_ = io.Discard // suppress unused import
	fmt.Println("All tests passed")
}

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

type constraintD struct{ u, v, x int }

// reference: lexicographically smallest array satisfying all OR constraints.
// Greedy left-to-right: set a[i] = lower[i], then propagate unmet requirements
// to later neighbours.
func referenceD(n int, cs []constraintD) []int {
	upper := make([]int, n+1)
	mask := (1 << 30) - 1
	for i := range upper {
		upper[i] = mask
	}
	for _, c := range cs {
		upper[c.u] &= c.x
		upper[c.v] &= c.x
	}

	lower := make([]int, n+1)
	type edge struct{ to, x int }
	adj := make([][]edge, n+1)

	for _, c := range cs {
		u, v, x := c.u, c.v, c.x
		if u == v {
			lower[u] |= x
		} else {
			lower[u] |= x &^ upper[v]
			lower[v] |= x &^ upper[u]
			if u < v {
				adj[u] = append(adj[u], edge{v, x})
			} else {
				adj[v] = append(adj[v], edge{u, x})
			}
		}
	}

	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = lower[i]
		for _, e := range adj[i] {
			lower[e.to] |= e.x &^ a[i]
		}
	}
	return a[1:]
}

// genCaseD generates a valid test case by picking a random array first,
// then deriving constraints from actual OR values — guaranteeing satisfiability.
func genCaseD(rng *rand.Rand) (n int, cs []constraintD) {
	n = rng.Intn(6) + 2
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(16)
	}
	q := rng.Intn(8) + 1
	cs = make([]constraintD, q)
	for k := 0; k < q; k++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		cs[k] = constraintD{u, v, a[u] | a[v]}
	}
	return
}

func runCaseD(bin string, n int, cs []constraintD) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(cs))
	for _, c := range cs {
		fmt.Fprintf(&sb, "%d %d %d\n", c.u, c.v, c.x)
	}
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	got := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		got[i] = v
	}

	// Verify all constraints are satisfied.
	for _, c := range cs {
		if got[c.u-1]|got[c.v-1] != c.x {
			return fmt.Errorf("constraint (%d,%d,%d) violated: a[%d]=%d, a[%d]=%d, OR=%d",
				c.u, c.v, c.x, c.u, got[c.u-1], c.v, got[c.v-1], got[c.u-1]|got[c.v-1])
		}
	}

	// Verify lexicographic minimality against reference.
	exp := referenceD(n, cs)
	for i := range exp {
		if got[i] != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 200; t++ {
		n, cs := genCaseD(rng)
		if err := runCaseD(bin, n, cs); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			fmt.Fprintf(os.Stderr, "%d %d\n", n, len(cs))
			for _, c := range cs {
				fmt.Fprintf(os.Stderr, "%d %d %d\n", c.u, c.v, c.x)
			}
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

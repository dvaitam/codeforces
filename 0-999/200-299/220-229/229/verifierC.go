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

type edgeC struct{ a, b int }

type testCaseC struct {
	n     int
	edges []edgeC
	ans   int64
}

func solveC(tc testCaseC) int64 {
	deg := make([]int64, tc.n)
	for _, e := range tc.edges {
		deg[e.a-1]++
		deg[e.b-1]++
	}
	N := int64(tc.n)
	if tc.n < 3 {
		return 0
	}
	total := N * (N - 1) * (N - 2) / 6
	var S int64
	for i := 0; i < tc.n; i++ {
		S += deg[i] * (N - 1 - deg[i])
	}
	return total - S/2
}

func genCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(6) + 1
	edges := make([]edgeC, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rng.Intn(2) == 0 {
				continue
			}
			edges = append(edges, edgeC{i, j})
		}
	}
	tc := testCaseC{n: n, edges: edges}
	tc.ans = solveC(tc)
	return tc
}

func runCaseC(bin string, tc testCaseC) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.a, e.b)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

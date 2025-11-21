package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	input     string
	n         int
	m         int
	balances  []int64
	refOutput string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	if out == "" {
		return fmt.Errorf("output empty")
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var mOut int
	if _, err := fmt.Fscan(reader, &mOut); err != nil {
		return fmt.Errorf("failed to read m': %v", err)
	}
	if mOut < 0 || mOut > 300000 {
		return fmt.Errorf("m' %d out of range", mOut)
	}
	edgesOut := make([]edge, mOut)
	seen := make(map[[2]int]bool)
	var total int64
	for i := 0; i < mOut; i++ {
		var u, v int
		var w int64
		if _, err := fmt.Fscan(reader, &u, &v, &w); err != nil {
			return fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
		if u == v {
			return fmt.Errorf("self-loop edge %d %d %d not allowed", u, v, w)
		}
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return fmt.Errorf("edge %d has invalid vertices %d -> %d", i+1, u, v)
		}
		if w <= 0 {
			return fmt.Errorf("edge %d weight must be positive", i+1)
		}
		key := [2]int{u, v}
		if seen[key] {
			return fmt.Errorf("duplicate edge between %d and %d", u, v)
		}
		seen[key] = true
		edgesOut[i] = edge{u, v, w}
		total += w
	}
	if err := validateNetBalances(edgesOut, tc.balances); err != nil {
		return err
	}
	refTotal := tc.refOutput
	wantTotal, _ := strconv.ParseInt(refTotal, 10, 64)
	if total != wantTotal {
		return fmt.Errorf("total debt %d expected %d", total, wantTotal)
	}
	return nil
}

func validateNetBalances(edges []edge, balances []int64) error {
	newBalance := make([]int64, len(balances))
	copy(newBalance, balances)
	for _, e := range edges {
		newBalance[e.u] -= e.w
		newBalance[e.v] += e.w
	}
	for i := 1; i < len(newBalance); i++ {
		if newBalance[i] != 0 {
			return fmt.Errorf("net balance mismatch at node %d", i)
		}
	}
	return nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest(1, []edge{}),
		makeTest(2, []edge{{u: 1, v: 2, w: 5}}),
		makeTest(3, []edge{
			{u: 1, v: 2, w: 10},
			{u: 2, v: 3, w: 5},
			{u: 1, v: 3, w: 5},
		}),
	}
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest())
	}
	return tests
}

func randomTest() testCase {
	n := rand.Intn(8) + 1
	m := rand.Intn(15)
	edges := make([]edge, 0, m)
	for i := 0; i < m; i++ {
		u := rand.Intn(n) + 1
		v := rand.Intn(n-1) + 1
		if v >= u {
			v++
		}
		w := int64(rand.Intn(1000) + 1)
		edges = append(edges, edge{u, v, w})
	}
	return makeTest(n, edges)
}

func makeTest(n int, edges []edge) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	if len(edges) > 0 {
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
		}
	}
	input := sb.String()
	bal := computeBalances(n, edges)
	ref := computeRef(n, edges)
	return testCase{
		input:     input,
		n:         n,
		m:         len(edges),
		balances:  bal,
		refOutput: fmt.Sprintf("%d", ref),
	}
}

func computeBalances(n int, edges []edge) []int64 {
	bal := make([]int64, n+1)
	for _, e := range edges {
		bal[e.u] -= e.w
		bal[e.v] += e.w
	}
	return bal
}

func computeRef(n int, edges []edge) int64 {
	balances := computeBalances(n, edges)
	var total int64
	for i := 1; i < len(balances); i++ {
		if balances[i] > 0 {
			total += balances[i]
		}
	}
	return total
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n      int
	parent []int
	c      []int
	exp    string
}

func dfsBuild(v int, children [][]int, c []int) ([]int, bool) {
	order := make([]int, 0)
	for _, ch := range children[v] {
		sub, ok := dfsBuild(ch, children, c)
		if !ok {
			return nil, false
		}
		order = append(order, sub...)
	}
	if c[v] > len(order) {
		return nil, false
	}
	idx := c[v]
	order = append(order, 0)
	copy(order[idx+1:], order[idx:])
	order[idx] = v
	return order, true
}

func solveB(n int, parent []int, c []int) string {
	children := make([][]int, n+1)
	root := 0
	for i := 1; i <= n; i++ {
		p := parent[i]
		if p == 0 {
			root = i
		} else {
			children[p] = append(children[p], i)
		}
	}
	order, ok := dfsBuild(root, children, c)
	if !ok || len(order) != n {
		return "NO"
	}
	ans := make([]int, n+1)
	for i, v := range order {
		ans[v] = i + 1
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(ans[i]))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
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

func genValidCase(rng *rand.Rand, n int) (parent []int, c []int) {
	parent = make([]int, n+1)
	parent[1] = 0
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	// assign random permutation to compute c
	perm := rng.Perm(n)
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = perm[i-1] + 1
	}
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parent[i]
		children[p] = append(children[p], i)
	}
	c = make([]int, n+1)
	var dfs func(int) []int
	dfs = func(v int) []int {
		order := make([]int, 0)
		for _, ch := range children[v] {
			sub := dfs(ch)
			order = append(order, sub...)
		}
		cnt := 0
		for _, u := range order {
			if val[u] < val[v] {
				cnt++
			}
		}
		c[v] = cnt
		order = append(order, v)
		return order
	}
	dfs(1)
	return parent, c
}

func generateTests() []testCaseB {
	rng := rand.New(rand.NewSource(2))
	cases := make([]testCaseB, 100)
	for i := range cases {
		n := rng.Intn(10) + 1
		if rng.Intn(2) == 0 {
			parent, cVals := genValidCase(rng, n)
			cases[i] = testCaseB{n: n, parent: parent, c: cVals, exp: solveB(n, parent, cVals)}
		} else {
			parent, cVals := genValidCase(rng, n)
			// introduce invalidity by increasing some c
			for j := 1; j <= n; j++ {
				if rng.Float64() < 0.3 {
					cVals[j] += rng.Intn(3) + n
				}
			}
			cases[i] = testCaseB{n: n, parent: parent, c: cVals, exp: solveB(n, parent, cVals)}
		}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j := 1; j <= tc.n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", tc.parent[j], tc.c[j])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

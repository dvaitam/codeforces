package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	input    string
	expected string
}

func can(g [][]int, n, k, x int) bool {
	sizes := make([]int, 0)
	var dfs func(v, p int) int
	dfs = func(v, p int) int {
		sum := 1
		for _, to := range g[v] {
			if to == p {
				continue
			}
			sum += dfs(to, v)
		}
		if sum >= x {
			sizes = append(sizes, sum)
			return 0
		}
		return sum
	}
	dfs(0, -1)
	if len(sizes) < k {
		return false
	}
	sort.Ints(sizes)
	total := 0
	for i := 0; i < k; i++ {
		total += sizes[i]
	}
	return n-total >= x
}

func solveOne(n, k int, edges [][2]int) int {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	low, high, ans := 1, n, 1
	for low <= high {
		mid := (low + high) / 2
		if can(g, n, k, mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return ans
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			edges[i] = [2]int{u, v}
		}
		ans := solveOne(n, k, edges)
		out.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(out.String())
}

func generateTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(48))
	var tests []test
	fixed := []string{
		"1\n1 0\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		k := rng.Intn(n)
		edges := generateTree(rng, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

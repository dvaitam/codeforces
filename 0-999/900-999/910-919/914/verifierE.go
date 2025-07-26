package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n     int
	edges [][2]int
	label []byte
}

func generateTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, testCase) {
	n := rng.Intn(5) + 2
	edges := generateTree(rng, n)
	labels := make([]byte, n)
	for i := 0; i < n; i++ {
		labels[i] = byte('a' + rng.Intn(3))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(string(labels))
	sb.WriteByte('\n')
	return sb.String(), testCase{n: n, edges: edges, label: labels}
}

func findPath(adj [][]int, u, v int) []int {
	n := len(adj) - 1
	prev := make([]int, n+1)
	for i := range prev {
		prev[i] = -1
	}
	q := []int{u}
	prev[u] = 0
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			break
		}
		for _, to := range adj[x] {
			if prev[to] == -1 {
				prev[to] = x
				q = append(q, to)
			}
		}
	}
	path := []int{}
	cur := v
	for cur != 0 {
		path = append(path, cur)
		if cur == u {
			break
		}
		cur = prev[cur]
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func palMask(path []int, labels []byte) int {
	mask := 0
	for _, v := range path {
		mask ^= 1 << (labels[v-1] - 'a')
	}
	return mask
}

func expected(tc testCase) []int {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	res := make([]int, n+1)
	for i := 1; i <= n; i++ {
		res[i] = 1 // path of length 0
	}
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			p := findPath(adj, u, v)
			if bits := palMask(p, tc.label); bits&(bits-1) == 0 {
				for _, x := range p {
					res[x]++
				}
			}
		}
	}
	return res[1:]
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCase(rng)
		expSlice := expected(tc)
		expStrings := make([]string, len(expSlice))
		for j, v := range expSlice {
			expStrings[j] = fmt.Sprintf("%d", v)
		}
		exp := strings.Join(expStrings, " ")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

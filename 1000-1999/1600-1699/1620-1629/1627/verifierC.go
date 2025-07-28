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
	"time"
)

type edge struct {
	to, id int
}

func solveC(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(r, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(r, &n)
		adj := make([][]edge, n+1)
		deg := make([]int, n+1)
		for i := 1; i < n; i++ {
			var u, v int
			fmt.Fscan(r, &u, &v)
			adj[u] = append(adj[u], edge{v, i})
			adj[v] = append(adj[v], edge{u, i})
			deg[u]++
			deg[v]++
		}
		ok := true
		for i := 1; i <= n; i++ {
			if deg[i] > 2 {
				ok = false
				break
			}
		}
		if !ok {
			out.WriteString("-1\n")
			continue
		}
		st := 1
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				st = i
				break
			}
		}
		ans := make([]int, n)
		prev, curr, parity := 0, st, 0
		for {
			found := false
			for _, e := range adj[curr] {
				if e.to == prev {
					continue
				}
				if parity == 0 {
					ans[e.id] = 2
				} else {
					ans[e.id] = 3
				}
				parity ^= 1
				prev = curr
				curr = e.to
				found = true
				break
			}
			if !found {
				break
			}
		}
		for i := 1; i < n; i++ {
			if i > 1 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.Itoa(ans[i]))
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func runProg(bin, input string) (string, error) {
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

func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(3))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 2
		edges := genTree(rng, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveC(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

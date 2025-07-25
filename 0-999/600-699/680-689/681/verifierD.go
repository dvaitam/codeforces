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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct{ u, v int }

func solve(n int, edges []edge, gifts []int) string {
	children := make([][]int, n+1)
	parent := make([]int, n+1)
	for _, e := range edges {
		children[e.u] = append(children[e.u], e.v)
		parent[e.v] = e.u
	}
	listVal := make([]bool, n+1)
	q := make([]int, n)
	head, tail := 0, 0
	for i := 1; i <= n; i++ {
		if gifts[i] >= 1 && gifts[i] <= n {
			listVal[gifts[i]] = true
		}
		if parent[i] == 0 {
			q[tail] = i
			tail++
			for head < tail {
				u := q[head]
				head++
				for _, v := range children[u] {
					q[tail] = v
					tail++
				}
			}
		}
	}
	var ans []int
	for i := tail - 1; i >= 0; i-- {
		u := q[i]
		if gifts[u] != u {
			p := parent[u]
			if p == 0 || gifts[p] != gifts[u] {
				return "-1"
			}
		}
		if listVal[u] {
			ans = append(ans, u)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", len(ans)))
	for _, x := range ans {
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	gifts := make([]int, n+1)
	for i := 1; i <= n; i++ {
		gifts[i] = rng.Intn(n + 2)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d", gifts[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String(), solve(n, edges, gifts)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

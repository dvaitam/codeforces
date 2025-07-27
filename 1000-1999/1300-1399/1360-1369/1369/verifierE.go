package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct {
	id int
	to int
}

func solveE(n, m int, cnt []int, a0 []int, a1 []int) (string, []int) {
	adj := make([][]edge, n+1)
	for i := 1; i <= m; i++ {
		adj[a0[i]] = append(adj[a0[i]], edge{i, a1[i]})
		adj[a1[i]] = append(adj[a1[i]], edge{i, a0[i]})
	}
	siz := make([]int, n+1)
	vis := make([]bool, n+1)
	inEdge := make([]bool, m+1)
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		siz[i] = len(adj[i])
		if cnt[i] >= siz[i] {
			vis[i] = true
			queue = append(queue, i)
		}
	}
	var stack []int
	qi := 0
	for qi < len(queue) {
		d := queue[qi]
		qi++
		for _, e := range adj[d] {
			if !inEdge[e.id] {
				inEdge[e.id] = true
				stack = append(stack, e.id)
			}
			if vis[e.to] {
				continue
			}
			siz[e.to]--
			if cnt[e.to] >= siz[e.to] {
				vis[e.to] = true
				queue = append(queue, e.to)
			}
		}
	}
	if len(stack) != m {
		return "DEAD", nil
	}
	order := make([]int, m)
	for i := m - 1; i >= 0; i-- {
		order[m-1-i] = stack[i]
	}
	return "ALIVE", order
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(6) + 2
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges) + 1
		cnt := make([]int, n+1)
		for i := 1; i <= n; i++ {
			cnt[i] = rand.Intn(3)
		}
		a0 := make([]int, m+1)
		a1 := make([]int, m+1)
		for i := 1; i <= m; i++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			for y == x {
				y = rand.Intn(n) + 1
			}
			a0[i] = x
			a1[i] = y
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, cnt[i])
		}
		sb.WriteByte('\n')
		for i := 1; i <= m; i++ {
			fmt.Fprintf(&sb, "%d %d\n", a0[i], a1[i])
		}
		expStatus, expOrder := solveE(n, m, append([]int(nil), cnt...), append([]int(nil), a0...), append([]int(nil), a1...))
		var expected strings.Builder
		expected.WriteString(expStatus)
		expected.WriteByte('\n')
		if expStatus == "ALIVE" {
			for i, v := range expOrder {
				if i > 0 {
					expected.WriteByte(' ')
				}
				fmt.Fprint(&expected, v)
			}
			expected.WriteByte('\n')
		}
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected.String()) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), expected.String(), output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveI(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
	}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[1] = 0
	q = append(q, 1)
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	g := 0
	for u := 1; u <= n; u++ {
		if dist[u] < 0 {
			continue
		}
		for _, v := range adj[u] {
			if dist[v] >= 0 {
				d := dist[u] + 1 - dist[v]
				if d < 0 {
					d = -d
				}
				g = gcd(g, d)
			}
		}
	}
	t := g
	if t <= 0 {
		t = 1
	}
	cams := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if dist[i] >= 0 && dist[i]%t == 0 {
			cams = append(cams, i)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	sb.WriteString(fmt.Sprintf("%d\n", len(cams)))
	for i, v := range cams {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCaseI(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges) + 1
	edges := make(map[[2]int]bool)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	count := 0
	for count < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if edges[key] {
			continue
		}
		edges[key] = true
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
		count++
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseI(rng)
	}
	for i, tc := range cases {
		expect := solveI(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

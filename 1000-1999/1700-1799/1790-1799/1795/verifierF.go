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

func expected(n int, edges [][2]int, k int, chips []int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dist := make([]int, n+1)
	owner := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	for i, v := range chips {
		idx := i + 1
		dist[v] = 0
		owner[v] = idx
		q = append(q, v)
	}
	head := 0
	for head < len(q) {
		v := q[head]
		head++
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				owner[to] = owner[v]
				q = append(q, to)
			} else if dist[to] == dist[v]+1 && owner[to] > owner[v] {
				owner[to] = owner[v]
			}
		}
	}
	maxDist := make([]int, k+1)
	for i := 1; i <= n; i++ {
		id := owner[i]
		if id >= 1 && dist[i] > maxDist[id] {
			maxDist[id] = dist[i]
		}
	}
	r := maxDist[1]
	j := 1
	for i := 1; i <= k; i++ {
		if maxDist[i] < r {
			r = maxDist[i]
			j = i
		}
	}
	moves := r*k + j - 1
	if moves > n-k {
		moves = n - k
	}
	return fmt.Sprintf("%d", moves)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 1
		edges := make([][2]int, n-1)
		perm := rand.Perm(n)
		for i := 1; i < n; i++ {
			u := perm[rng.Intn(i)] + 1
			v := perm[i] + 1
			edges[i-1] = [2]int{u, v}
		}
		k := rng.Intn(n) + 1
		chips := rand.Perm(n)[:k]
		for i := 0; i < k; i++ {
			chips[i]++
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for i, c := range chips {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOut := expected(n, edges, k, chips)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

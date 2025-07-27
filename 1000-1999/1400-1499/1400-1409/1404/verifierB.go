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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func bfsDist(src, tgt int, adj [][]int, n int) int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[src] = 0
	q = append(q, src)
	for i := 0; i < len(q); i++ {
		u := q[i]
		if u == tgt {
			break
		}
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist[tgt]
}

func bfsFarthest(start int, adj [][]int, n int) (node, dist int) {
	d := make([]int, n)
	for i := range d {
		d[i] = -1
	}
	q := make([]int, 0, n)
	d[start] = 0
	q = append(q, start)
	node = start
	dist = 0
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if d[v] == -1 {
				d[v] = d[u] + 1
				q = append(q, v)
				if d[v] > dist {
					dist = d[v]
					node = v
				}
			}
		}
	}
	return
}

func solve(reader *bufio.Reader) string {
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; t > 0; t-- {
		var n, a, b, da, db int
		fmt.Fscan(reader, &n, &a, &b, &da, &db)
		a--
		b--
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		distAB := bfsDist(a, b, adj, n)
		if distAB <= da {
			sb.WriteString("Alice\n")
			continue
		}
		far, _ := bfsFarthest(0, adj, n)
		_, diam := bfsFarthest(far, adj, n)
		if db <= 2*da || diam <= 2*da {
			sb.WriteString("Alice\n")
		} else {
			sb.WriteString("Bob\n")
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for ; t > 0; t-- {
		n := rng.Intn(10) + 2
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		da := rng.Intn(n) + 1
		db := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, a, b, da, db)
		// generate tree
		parents := make([]int, n)
		for i := 1; i < n; i++ {
			p := rng.Intn(i)
			parents[i] = p
			fmt.Fprintf(&sb, "%d %d\n", i+1, p+1)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Embedded reference solver for 576D

type edgeD struct {
	a, b int
	d    int64
}

type bitsD [3]uint64

func setBitD(b *bitsD, i int) {
	b[i>>6] |= 1 << (uint(i) & 63)
}

func isEmptyD(b bitsD) bool {
	return b[0] == 0 && b[1] == 0 && b[2] == 0
}

func addEdgeD(rows []bitsD, adj, rev [][]int, u, v int) {
	if ((rows[u][v>>6] >> (uint(v) & 63)) & 1) != 0 {
		return
	}
	rows[u][v>>6] |= 1 << (uint(v) & 63)
	adj[u] = append(adj[u], v)
	rev[v] = append(rev[v], u)
}

func shortestToTargetD(s bitsD, rev [][]int, target, n int) int {
	if isEmptyD(s) {
		return -1
	}
	if ((s[target>>6] >> (uint(target) & 63)) & 1) != 0 {
		return 0
	}
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := make([]int, n)
	head, tail := 0, 0
	dist[target] = 0
	q[tail] = target
	tail++
	for head < tail {
		v := q[head]
		head++
		for _, u := range rev[v] {
			if dist[u] == -1 {
				dist[u] = dist[v] + 1
				q[tail] = u
				tail++
			}
		}
	}
	best := -1
	for w := 0; w < 3; w++ {
		x := s[w]
		for x != 0 {
			t := bits.TrailingZeros64(x)
			i := 64*w + t
			if i < n && dist[i] != -1 && (best == -1 || dist[i] < best) {
				best = dist[i]
			}
			x &= x - 1
		}
	}
	return best
}

func squareD(prev []bitsD, n int) []bitsD {
	next := make([]bitsD, n)
	for i := 0; i < n; i++ {
		var row bitsD
		for w := 0; w < 3; w++ {
			x := prev[i][w]
			for x != 0 {
				t := bits.TrailingZeros64(x)
				j := 64*w + t
				if j < n {
					row[0] |= prev[j][0]
					row[1] |= prev[j][1]
					row[2] |= prev[j][2]
				}
				x &= x - 1
			}
		}
		next[i] = row
	}
	return next
}

func applyD(vec bitsD, mat []bitsD, n int) bitsD {
	var res bitsD
	for w := 0; w < 3; w++ {
		x := vec[w]
		for x != 0 {
			t := bits.TrailingZeros64(x)
			i := 64*w + t
			if i < n {
				res[0] |= mat[i][0]
				res[1] |= mat[i][1]
				res[2] |= mat[i][2]
			}
			x &= x - 1
		}
	}
	return res
}

func advanceExactD(s bitsD, rows []bitsD, n int, steps int64) bitsD {
	if steps == 0 {
		return s
	}
	maxP := 0
	for x := steps; x > 0; x >>= 1 {
		maxP++
	}
	pows := make([][]bitsD, maxP)
	pows[0] = make([]bitsD, n)
	copy(pows[0], rows)
	for i := 1; i < maxP; i++ {
		pows[i] = squareD(pows[i-1], n)
	}
	cur := s
	bit := 0
	for steps > 0 {
		if (steps & 1) != 0 {
			cur = applyD(cur, pows[bit], n)
			if isEmptyD(cur) {
				return cur
			}
		}
		steps >>= 1
		bit++
	}
	return cur
}

func solveD(n, m int, edges []edgeD) string {
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].d != edges[j].d {
			return edges[i].d < edges[j].d
		}
		if edges[i].a != edges[j].a {
			return edges[i].a < edges[j].a
		}
		return edges[i].b < edges[j].b
	})

	rows := make([]bitsD, n)
	adj := make([][]int, n)
	rev := make([][]int, n)

	idx := 0
	for idx < m && edges[idx].d == 0 {
		addEdgeD(rows, adj, rev, edges[idx].a, edges[idx].b)
		idx++
	}

	var cur bitsD
	setBitD(&cur, 0)
	var flights int64

	for idx < m {
		nextDD := edges[idx].d
		gap := nextDD - flights

		dist := shortestToTargetD(cur, rev, n-1, n)
		if dist != -1 && int64(dist) <= gap {
			return fmt.Sprint(flights + int64(dist))
		}

		cur = advanceExactD(cur, rows, n, gap)
		if isEmptyD(cur) {
			return "Impossible"
		}

		flights = nextDD
		for idx < m && edges[idx].d == flights {
			addEdgeD(rows, adj, rev, edges[idx].a, edges[idx].b)
			idx++
		}
	}

	dist := shortestToTargetD(cur, rev, n-1, n)
	if dist == -1 {
		return "Impossible"
	}
	return fmt.Sprint(flights + int64(dist))
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stderr = new(bytes.Buffer)
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rand.Seed(4)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 2
		m := rand.Intn(8) + 1
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < m; j++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			d := rand.Intn(20)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, d))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	if !filepath.IsAbs(cand) {
		abs, err := filepath.Abs(cand)
		if err == nil {
			cand = abs
		}
	}

	tests := genTests()
	for idx, input := range tests {
		// Parse input to get n, m, edges
		sc := strings.Fields(input)
		pos := 0
		atoi := func() int {
			v := 0
			for _, c := range sc[pos] {
				v = v*10 + int(c-'0')
			}
			pos++
			return v
		}
		n := atoi()
		m := atoi()
		edges := make([]edgeD, m)
		for j := 0; j < m; j++ {
			a := atoi()
			b := atoi()
			d := atoi()
			edges[j] = edgeD{a - 1, b - 1, int64(d)}
		}

		exp := solveD(n, m, edges)
		out, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("runtime error on test %d\n", idx+1)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

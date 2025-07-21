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

func computeLR(p []int) ([]int, []int) {
	n := len(p)
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i
		for j := i - 1; j >= 0; j-- {
			if p[j] > p[i] {
				l[i] = j
				break
			}
		}
		r[i] = i
		for j := i + 1; j < n; j++ {
			if p[j] > p[i] {
				r[i] = j
				break
			}
		}
	}
	return l, r
}

func diameter(n int, edges [][2]int) int {
	g := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	bfs := func(start int) (int, int) {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		idx := 0
		for idx < len(q) {
			v := q[idx]
			idx++
			for _, u := range g[v] {
				if dist[u] == -1 {
					dist[u] = dist[v] + 1
					q = append(q, u)
				}
			}
		}
		far := start
		for i, d := range dist {
			if d > dist[far] {
				far = i
			}
			if d == -1 {
				return -1, -1
			}
		}
		return far, dist[far]
	}
	f, _ := bfs(0)
	if f == -1 {
		return -1
	}
	_, d := bfs(f)
	return d
}

func brute(n int, p []int, s string) int {
	l, r := computeLR(p)
	var qIdx []int
	for i, ch := range s {
		if ch == '?' {
			qIdx = append(qIdx, i)
		}
	}
	best := -1
	total := 1 << len(qIdx)
	for mask := 0; mask < total; mask++ {
		edges := make([][2]int, 0, n)
		connected := true
		for i := 0; i < n; i++ {
			c := s[i]
			if c == '?' {
				bit := 0
				for j, idx := range qIdx {
					if idx == i {
						bit = (mask >> j) & 1
						break
					}
				}
				if bit == 0 {
					c = 'L'
				} else {
					c = 'R'
				}
			}
			var to int
			if c == 'L' {
				to = l[i]
			} else {
				to = r[i]
			}
			if to == i {
				connected = false
				break
			}
			edges = append(edges, [2]int{i, to})
		}
		if !connected {
			continue
		}
		d := diameter(n, edges)
		if d > best {
			best = d
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 2
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	s := make([]byte, n)
	chars := []byte{'L', 'R', '?'}
	for i := range s {
		s[i] = chars[rng.Intn(len(chars))]
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(string(s))
	sb.WriteByte('\n')
	cp := append([]int(nil), p...)
	exp := brute(n, cp, string(s))
	return sb.String(), exp
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

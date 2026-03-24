package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	to, rev         int
	cap, flow, cost int
}

const SOLINF int = 1e9

func solve(n int, s string, b []int) int {
	S := 0
	T := 26 + n/2 + 1
	numNodes := T + 1

	graph := make([][]Edge, numNodes)

	addEdge := func(u, v, cap, cost int) {
		graph[u] = append(graph[u], Edge{v, len(graph[v]), cap, 0, cost})
		graph[v] = append(graph[v], Edge{u, len(graph[u]) - 1, 0, 0, -cost})
	}

	cnt := make([]int, 26)
	for i := 0; i < n; i++ {
		cnt[s[i]-'a']++
	}

	for c := 0; c < 26; c++ {
		if cnt[c] > 0 {
			addEdge(S, c+1, cnt[c], 0)
		}
	}

	for p := 1; p <= n/2; p++ {
		i := p - 1
		j := n - p
		charI := int(s[i] - 'a')
		charJ := int(s[j] - 'a')
		bI := b[i]
		bJ := b[j]

		pairNode := 26 + p
		addEdge(pairNode, T, 2, 0)

		for c := 0; c < 26; c++ {
			profit := 0
			if charI != charJ {
				if c == charI {
					profit = bI
				} else if c == charJ {
					profit = bJ
				}
			} else {
				if c == charI {
					if bI > bJ {
						profit = bI
					} else {
						profit = bJ
					}
				}
			}
			addEdge(c+1, pairNode, 1, -profit)
		}
	}

	totalFlow := 0
	minCost := 0

	for totalFlow < n {
		dist := make([]int, numNodes)
		for i := range dist {
			dist[i] = SOLINF
		}
		parentEdge := make([]int, numNodes)
		parentNode := make([]int, numNodes)
		inQueue := make([]bool, numNodes)

		dist[S] = 0
		queue := []int{S}
		inQueue[S] = true

		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			inQueue[u] = false

			for i, e := range graph[u] {
				if e.cap-e.flow > 0 && dist[e.to] > dist[u]+e.cost {
					dist[e.to] = dist[u] + e.cost
					parentNode[e.to] = u
					parentEdge[e.to] = i
					if !inQueue[e.to] {
						queue = append(queue, e.to)
						inQueue[e.to] = true
					}
				}
			}
		}

		if dist[T] == SOLINF {
			break
		}

		push := SOLINF
		curr := T
		for curr != S {
			p := parentNode[curr]
			idx := parentEdge[curr]
			rem := graph[p][idx].cap - graph[p][idx].flow
			if rem < push {
				push = rem
			}
			curr = p
		}

		totalFlow += push
		minCost += push * dist[T]

		curr = T
		for curr != S {
			p := parentNode[curr]
			idx := parentEdge[curr]
			revIdx := graph[p][idx].rev
			graph[p][idx].flow += push
			graph[curr][revIdx].flow -= push
			curr = p
		}
	}

	return -minCost
}

func generateString(n int) (string, string) {
	t := make([]byte, n)
	for i := 0; i < n/2; i++ {
		a := byte('a' + rand.Intn(26))
		b := byte('a' + rand.Intn(26))
		for b == a {
			b = byte('a' + rand.Intn(26))
		}
		t[i] = a
		t[n-1-i] = b
	}
	perm := rand.Perm(n)
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = t[perm[i]]
	}
	return string(s), string(t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(10)*2 + 2 // even >=2
		s, _ := generateString(n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = rand.Intn(100) + 1
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		fmt.Fprintln(&input, s)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d ", b[i])
		}
		fmt.Fprintln(&input)
		expected := solve(n, s, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(expected) {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Edmonds' blossom algorithm for general graph maximum matching.
type MaximumMatching struct {
	n                     int
	mate, parity          []int
	blossom, parent, last []int
	adj                   [][]int
	bfsQueue              []int
	size, timer           int
}

func NewMaximumMatching(n int) *MaximumMatching {
	mm := &MaximumMatching{n: n}
	mm.mate = make([]int, n)
	mm.parity = make([]int, n)
	mm.blossom = make([]int, n)
	mm.parent = make([]int, n)
	mm.last = make([]int, n)
	mm.adj = make([][]int, n)
	for i := 0; i < n; i++ {
		mm.mate[i] = -1
	}
	return mm
}

func (mm *MaximumMatching) AddEdge(u, v int) {
	mm.adj[u] = append(mm.adj[u], v)
	mm.adj[v] = append(mm.adj[v], u)
}

func (mm *MaximumMatching) augment(u int) {
	for u != -1 {
		v := mm.parent[u]
		w := mm.mate[v]
		mm.mate[u] = v
		mm.mate[v] = u
		u = w
	}
}

func (mm *MaximumMatching) lca(u, v int) int {
	mm.timer++
	for {
		if u != -1 {
			if mm.last[u] == mm.timer {
				return u
			}
			mm.last[u] = mm.timer
			if mm.mate[u] == -1 {
				u = -1
			} else {
				u = mm.blossom[mm.parent[mm.mate[u]]]
			}
		}
		u, v = v, u
	}
}

func (mm *MaximumMatching) merge(u, v, p int) {
	for mm.blossom[u] != p {
		mm.parent[u] = v
		v = mm.mate[u]
		if mm.parity[v] == 1 {
			mm.parity[v] = 0
			mm.bfsQueue = append(mm.bfsQueue, v)
		}
		mm.blossom[u] = p
		mm.blossom[v] = p
		u = mm.parent[v]
	}
}

func (mm *MaximumMatching) bfs(root int) bool {
	for i := 0; i < mm.n; i++ {
		mm.parity[i] = -1
		mm.parent[i] = -1
		mm.blossom[i] = i
	}
	mm.bfsQueue = mm.bfsQueue[:0]
	mm.bfsQueue = append(mm.bfsQueue, root)
	mm.parity[root] = 0
	for qi := 0; qi < len(mm.bfsQueue); qi++ {
		u := mm.bfsQueue[qi]
		for _, v := range mm.adj[u] {
			if mm.parity[v] == -1 {
				mm.parity[v] = 1
				mm.parent[v] = u
				if mm.mate[v] == -1 {
					mm.augment(v)
					return true
				}
				mm.parity[mm.mate[v]] = 0
				mm.bfsQueue = append(mm.bfsQueue, mm.mate[v])
			} else if mm.parity[v] == 0 && mm.blossom[u] != mm.blossom[v] {
				p := mm.lca(mm.blossom[u], mm.blossom[v])
				mm.merge(u, v, p)
				mm.merge(v, u, p)
			}
		}
	}
	return false
}

func (mm *MaximumMatching) SolveLimit(limit int) int {
	for i := 0; i < mm.n; i++ {
		if mm.mate[i] == -1 {
			if mm.bfs(i) {
				mm.size++
				if limit > 0 && mm.size >= limit {
					break
				}
			}
		}
	}
	return mm.size
}

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) next() (string, bool) {
	var sb strings.Builder
	for {
		c, err := fs.r.ReadByte()
		if err != nil {
			if sb.Len() == 0 {
				return "", false
			}
			break
		}
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			if sb.Len() == 0 {
				continue
			}
			break
		}
		sb.WriteByte(c)
	}
	return sb.String(), true
}

func readInt(tok string) int {
	val, _ := strconv.Atoi(tok)
	return val
}

func main() {
	in := newScanner()
	first, ok := in.next()
	if !ok {
		return
	}

	var t int
	pending := ""

	if first == "manual" {
		tok, _ := in.next()
		t = readInt(tok)
	} else {
		t = readInt(first)
		if tok, has := in.next(); has {
			if tok != "manual" {
				pending = tok
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n int
		if pending != "" {
			n = readInt(pending)
			pending = ""
		} else {
			tok, _ := in.next()
			n = readInt(tok)
		}

		totalEdges := n * (n - 1) / 2
		var builder strings.Builder
		for builder.Len() < totalEdges {
			tok, ok := in.next()
			if !ok {
				break
			}
			builder.WriteString(tok)
		}
		s := builder.String()
		if len(s) < totalEdges {
			s += strings.Repeat("0", totalEdges-len(s))
		} else if len(s) > totalEdges {
			s = s[:totalEdges]
		}

		mmFriend := NewMaximumMatching(n)
		mmEnemy := NewMaximumMatching(n)

		pos := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if s[pos] == '1' {
					mmFriend.AddEdge(i, j)
				} else {
					mmEnemy.AddEdge(i, j)
				}
				pos++
			}
		}

		k := (n + 1) / 3
		friendsMatch := mmFriend.SolveLimit(k)
		var chosen *MaximumMatching
		if friendsMatch >= k {
			chosen = mmFriend
		} else {
			mmEnemy.SolveLimit(k)
			chosen = mmEnemy
		}

		pairs := make([][2]int, 0, k)
		for i := 0; i < n && len(pairs) < k; i++ {
			if chosen.mate[i] != -1 && i < chosen.mate[i] {
				pairs = append(pairs, [2]int{i + 1, chosen.mate[i] + 1})
			}
		}

		if len(pairs) < k {
			k = len(pairs)
		}

		fmt.Fprintln(out, k)
		for i, p := range pairs[:k] {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprintf(out, "%d %d", p[0], p[1])
		}
		fmt.Fprintln(out)
	}
}

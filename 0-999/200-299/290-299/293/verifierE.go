package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type edge struct{ to, w int }

type pair struct{ d, w int }

type bit struct {
	n    int
	tree []int
	ver  []int
	time int
}

func newBIT(n int) *bit {
	return &bit{n: n, tree: make([]int, n+1), ver: make([]int, n+1), time: 1}
}

func (b *bit) reset() { b.time++ }

func (b *bit) update(i, v int) {
	for ; i <= b.n; i += i & -i {
		if b.ver[i] != b.time {
			b.ver[i] = b.time
			b.tree[i] = 0
		}
		b.tree[i] += v
	}
}

func (b *bit) query(i int) int {
	if i > b.n {
		i = b.n
	}
	res := 0
	for ; i > 0; i -= i & -i {
		if b.ver[i] == b.time {
			res += b.tree[i]
		}
	}
	return res
}

func solve(input string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) < 3 {
		return "", fmt.Errorf("invalid input")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	L, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	W, err := strconv.Atoi(fields[2])
	if err != nil {
		return "", err
	}
	if len(fields) != 3+2*(n-1) {
		return "", fmt.Errorf("expected %d edge values, got %d", 2*(n-1), len(fields)-3)
	}

	adj := make([][]edge, n+1)
	pos := 3
	for i := 2; i <= n; i++ {
		p, _ := strconv.Atoi(fields[pos])
		w, _ := strconv.Atoi(fields[pos+1])
		pos += 2
		adj[i] = append(adj[i], edge{to: p, w: w})
		adj[p] = append(adj[p], edge{to: i, w: w})
	}

	removed := make([]bool, n+1)
	subSize := make([]int, n+1)
	b := newBIT(L + 1)
	var total int64

	var dfsSize func(u, p int) int
	dfsSize = func(u, p int) int {
		subSize[u] = 1
		for _, e := range adj[u] {
			if e.to != p && !removed[e.to] {
				subSize[u] += dfsSize(e.to, u)
			}
		}
		return subSize[u]
	}

	var findCentroid func(u, p, sz int) int
	findCentroid = func(u, p, sz int) int {
		for _, e := range adj[u] {
			if e.to != p && !removed[e.to] && subSize[e.to] > sz/2 {
				return findCentroid(e.to, u, sz)
			}
		}
		return u
	}

	collect := func(u, p, ew int) []pair {
		type st struct{ u, p, d, w int }
		stack := []st{{u, p, 1, ew}}
		var res []pair
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if cur.d > L || cur.w > W {
				continue
			}
			res = append(res, pair{d: cur.d, w: cur.w})
			for _, e := range adj[cur.u] {
				if e.to != cur.p && !removed[e.to] {
					stack = append(stack, st{u: e.to, p: cur.u, d: cur.d + 1, w: cur.w + e.w})
				}
			}
		}
		return res
	}

	var decompose func(int)
	decompose = func(u int) {
		sz := dfsSize(u, -1)
		c := findCentroid(u, -1, sz)
		removed[c] = true

		vec := []pair{{0, 0}}
		for _, e := range adj[c] {
			if removed[e.to] {
				continue
			}
			sub := collect(e.to, c, e.w)
			// sort by weight
			sort.Slice(sub, func(i, j int) bool { return sub[i].w < sub[j].w })
			sort.Slice(vec, func(i, j int) bool { return vec[i].w < vec[j].w })
			b.reset()
			ptr := 0
			for _, p := range sub {
				limW := W - p.w
				for ptr < len(vec) && vec[ptr].w <= limW {
					if vec[ptr].d <= L {
						b.update(vec[ptr].d+1, 1)
					}
					ptr++
				}
				remD := L - p.d
				if remD >= 0 {
					total += int64(b.query(remD + 1))
				}
			}
			vec = mergePairs(vec, sub)
		}
		for _, e := range adj[c] {
			if !removed[e.to] {
				decompose(e.to)
			}
		}
	}

	decompose(1)
	return fmt.Sprint(total), nil
}

func sortPairs(a []pair) {
	// simple insertion sort sufficient for small slices; replace with sort.Slice if desired
	for i := 1; i < len(a); i++ {
		j := i
		for j > 0 && a[j-1].w > a[j].w {
			a[j-1], a[j] = a[j], a[j-1]
			j--
		}
	}
}

func mergePairs(a, b []pair) []pair {
	res := make([]pair, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i].w < b[j].w {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	res = append(res, a[i:]...)
	res = append(res, b[j:]...)
	return res
}

var testcases = []string{
	"5 1 14 1 9 2 8 2 2 2 10",
	"1 1 15",
	"6 6 7 1 3 1 9 3 2 2 9 5 2",
	"1 1 8",
	"1 1 19",
	"4 3 18 1 0 1 2 2 0",
	"3 2 8 1 6 1 4",
	"6 2 6 1 9 1 3 3 3 4 5 1 4",
	"4 2 10 1 3 2 1 2 10",
	"6 3 7 1 4 1 5 1 7 1 6 4 0",
	"6 3 16 1 7 1 2 1 4 1 6 3 3",
	"4 4 8 1 2 2 4 1 0",
	"2 2 8 1 10",
	"1 1 9",
	"3 3 6 1 6 1 10",
	"1 1 11",
	"6 4 14 1 10 2 6 2 6 3 9 1 1",
	"5 4 4 1 0 2 10 2 3 4 10",
	"6 1 13 1 10 2 9 1 4 3 8 1 8",
	"6 6 10 1 2 1 1 2 5 4 3 2 5",
	"6 5 14 1 8 1 7 2 0 1 0 1 5",
	"6 2 9 1 8 2 1 1 8 3 0 2 1",
	"6 2 14 1 2 2 4 3 2 4 9 1 8",
	"4 1 3 1 4 1 10 2 3",
	"3 3 12 1 7 2 2",
	"4 1 18 1 9 1 4 2 8",
	"3 2 20 1 4 2 4",
	"2 1 7 1 5",
	"6 1 11 1 9 2 10 1 0 1 7 5 2",
	"2 2 19 1 8",
	"2 1 15 1 7",
	"4 3 0 1 9 1 1 2 2",
	"6 1 16 1 9 2 3 1 6 1 8 4 5",
	"6 1 6 1 6 2 8 3 4 1 0 1 5",
	"4 4 12 1 4 2 6 3 0",
	"6 4 20 1 10 2 2 2 10 4 6 5 9",
	"3 2 17 1 4 2 9",
	"6 2 11 1 6 2 7 2 1 2 5 2 6",
	"1 1 4",
	"2 2 2 1 2",
	"5 1 7 1 0 2 10 2 2 3 6",
	"4 3 17 1 2 2 2 1 10",
	"1 1 9",
	"6 3 9 1 0 2 5 2 1 1 10 3 3",
	"4 3 17 1 6 1 6 1 0",
	"1 1 0",
	"3 2 1 1 0 1 2",
	"5 2 15 1 9 1 0 3 8 2 1",
	"1 1 19",
	"6 1 6 1 7 2 1 1 7 4 8 2 8",
	"1 1 11",
	"6 4 16 1 1 1 8 1 1 4 5 5 5",
	"5 5 17 1 4 1 9 3 1 2 9",
	"4 3 9 1 8 2 0 1 0",
	"2 2 6 1 2",
	"1 1 15",
	"6 2 8 1 2 1 2 3 5 3 3 1 8",
	"4 2 20 1 5 1 10 1 9",
	"6 6 7 1 6 1 4 2 8 3 0 3 3",
	"4 1 13 1 2 2 0 3 4",
	"6 1 10 1 0 2 3 3 2 3 5 5 10",
	"1 1 15",
	"2 1 13 1 10",
	"1 1 0",
	"1 1 10",
	"2 2 7 1 0",
	"6 5 6 1 8 2 9 3 5 3 5 1 0",
	"4 3 5 1 3 2 9 1 10",
	"2 2 13 1 2",
	"2 2 12 1 10",
	"6 2 8 1 2 1 4 3 2 1 2 3 2",
	"4 3 16 1 5 2 7 2 1",
	"3 3 10 1 4 2 6",
	"5 5 0 1 10 1 6 2 10 4 0",
	"1 1 17",
	"1 1 0",
	"1 1 12",
	"5 3 17 1 8 2 8 3 7 2 9",
	"1 1 6",
	"3 2 2 1 1 1 9",
	"4 2 15 1 7 1 7 2 7",
	"3 3 16 1 0 2 8",
	"3 2 17 1 4 1 8",
	"2 1 2 1 3",
	"1 1 15",
	"4 1 18 1 2 2 9 1 9",
	"3 1 9 1 2 2 5",
	"1 1 3",
	"1 1 9",
	"6 3 18 1 2 1 10 2 6 4 10 2 1",
	"4 3 15 1 10 2 6 1 6",
	"2 2 20 1 4",
	"1 1 14",
	"3 2 9 1 5 1 2",
	"4 2 10 1 7 1 1 3 9",
	"5 4 10 1 1 2 2 2 5 1 1",
	"4 3 10 1 5 1 3 3 10",
	"3 3 5 1 8 1 7",
	"2 1 16 1 1",
	"1 1 2",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := strings.TrimSpace(tc) + "\n"

		expected, err := solve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}

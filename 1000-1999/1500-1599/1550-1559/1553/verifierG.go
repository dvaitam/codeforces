package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Uint64Slice []uint64

func (x Uint64Slice) Len() int           { return len(x) }
func (x Uint64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Embedded correct solver for 1553/G
func solveG(input string) string {
	idx := 0
	data := []byte(input)
	readInt := func() int {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		val := 0
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return val
	}

	n := readInt()
	q := readInt()

	maxA := 1000005
	spf := make([]int, maxA+1)
	for i := 2; i <= maxA; i++ {
		spf[i] = i
	}
	for i := 2; i*i <= maxA; i++ {
		if spf[i] == i {
			for j := i * i; j <= maxA; j += i {
				if spf[j] == j {
					spf[j] = i
				}
			}
		}
	}

	parent := make([]int, maxA+1)
	for i := 1; i <= maxA; i++ {
		parent[i] = i
	}

	var find func(int) int
	find = func(i int) int {
		root := i
		for parent[root] != root {
			root = parent[root]
		}
		curr := i
		for curr != root {
			nxt := parent[curr]
			parent[curr] = root
			curr = nxt
		}
		return root
	}

	union := func(i, j int) {
		rootI := find(i)
		rootJ := find(j)
		if rootI != rootJ {
			parent[rootI] = rootJ
		}
	}

	hasPrime := make([]bool, maxA+1)
	a := make([]int, n+1)

	for i := 1; i <= n; i++ {
		a[i] = readInt()
		v := a[i]
		for v > 1 {
			p := spf[v]
			hasPrime[p] = true
			union(a[i], p)
			for v%p == 0 {
				v /= p
			}
		}
	}

	pairs := make([]uint64, 0, n*21)

	for i := 1; i <= n; i++ {
		S := make([]int, 0, 8)
		S = append(S, find(a[i]))

		v := a[i] + 1
		for v > 1 {
			p := spf[v]
			if hasPrime[p] {
				S = append(S, find(p))
			}
			for v%p == 0 {
				v /= p
			}
		}

		for j := 0; j < len(S); j++ {
			for k := j + 1; k < len(S); k++ {
				u, v := S[j], S[k]
				if u == v {
					continue
				}
				if u > v {
					u, v = v, u
				}
				pairs = append(pairs, (uint64(u)<<32)|uint64(v))
			}
		}
	}

	sort.Sort(Uint64Slice(pairs))

	var uniquePairs []uint64
	if len(pairs) > 0 {
		uniquePairs = make([]uint64, 0, len(pairs))
		uniquePairs = append(uniquePairs, pairs[0])
		for i := 1; i < len(pairs); i++ {
			if pairs[i] != pairs[i-1] {
				uniquePairs = append(uniquePairs, pairs[i])
			}
		}
	}
	pairs = uniquePairs

	var out strings.Builder
	w := bufio.NewWriter(&out)

	for qIdx := 0; qIdx < q; qIdx++ {
		s := readInt()
		t := readInt()

		u := find(a[s])
		v := find(a[t])

		if u == v {
			fmt.Fprintln(w, 0)
			continue
		}

		if u > v {
			u, v = v, u
		}

		key := (uint64(u) << 32) | uint64(v)

		idx2 := sort.Search(len(pairs), func(i int) bool {
			return pairs[i] >= key
		})

		if idx2 < len(pairs) && pairs[idx2] == key {
			fmt.Fprintln(w, 1)
		} else {
			fmt.Fprintln(w, 2)
		}
	}
	w.Flush()

	return strings.TrimSpace(out.String())
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func randDistinct(rng *rand.Rand, n, lo, hi int) []int {
	pool := make([]int, 0, hi-lo+1)
	for v := lo; v <= hi; v++ {
		pool = append(pool, v)
	}
	rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	return pool[:n]
}

type testCase struct {
	input    string
	expected string
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(7))
	cases := []testCase{}

	// Fixed test
	fixed := []struct {
		n, q    int
		arr     []int
		queries [][2]int
	}{
		{2, 1, []int{2, 3}, [][2]int{{1, 2}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", f.n, f.q)
		for i, v := range f.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for _, qv := range f.queries {
			fmt.Fprintf(&sb, "%d %d\n", qv[0], qv[1])
		}
		inp := sb.String()
		exp := solveG(inp)
		cases = append(cases, testCase{inp, exp})
	}

	for len(cases) < 100 {
		n := rng.Intn(4) + 2
		q := rng.Intn(3) + 1
		arr := randDistinct(rng, n, 2, 20)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			s := rng.Intn(n) + 1
			t := rng.Intn(n-1) + 1
			if t >= s {
				t++
			}
			fmt.Fprintf(&sb, "%d %d\n", s, t)
		}
		inp := sb.String()
		exp := solveG(inp)
		cases = append(cases, testCase{inp, exp})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierG.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input  string
	output string
}

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func matchPattern(s, t string) bool {
	for i := 0; i < len(s); i++ {
		if t[i] != '_' && t[i] != s[i] {
			return false
		}
	}
	return true
}

func solveCase(N, M, K int, A []string, B []string, P []int) (bool, []int) {
	bMx := 1 << K
	nodes := 1
	for i := 0; i < K; i++ {
		nodes *= 27
	}
	isValid := make([]bool, nodes)
	id := make([]int, nodes)
	par := make([]int, nodes)
	adj := make([][]int, nodes)
	getCharId := func(ch byte) int {
		if ch == '_' {
			return 26
		}
		return int(ch - 'a')
	}
	hsh := func(s string) int {
		res := 0
		mult := 1
		for i := 0; i < K; i++ {
			res += mult * getCharId(s[i])
			mult *= 27
		}
		return res
	}
	for i := 0; i < N; i++ {
		u := hsh(A[i])
		isValid[u] = true
		id[u] = i + 1
	}
	for i := 0; i < M; i++ {
		p := P[i] - 1
		if !matchPattern(B[i], A[p]) {
			return false, nil
		}
		sBytes := []byte(B[i])
		t := A[p]
		u := hsh(t)
		for mask := 0; mask < bMx; mask++ {
			tmp := make([]byte, K)
			copy(tmp, sBytes)
			for j := 0; j < K; j++ {
				if mask&(1<<j) != 0 {
					tmp[j] = '_'
				}
			}
			v := hsh(string(tmp))
			if u == v || !isValid[u] || !isValid[v] {
				continue
			}
			adj[u] = append(adj[u], v)
			par[v]++
		}
	}
	h := &IntHeap{}
	heap.Init(h)
	for i := 0; i < N; i++ {
		u := hsh(A[i])
		if par[u] == 0 {
			heap.Push(h, u)
		}
	}
	result := make([]int, 0, N)
	for len(result) < N {
		if h.Len() == 0 {
			return false, nil
		}
		u := heap.Pop(h).(int)
		result = append(result, id[u])
		for _, v := range adj[u] {
			par[v]--
			if par[v] == 0 {
				heap.Push(h, v)
			}
		}
	}
	return true, result
}

func buildCase(N, M, K int, A []string, B []string, P []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", N, M, K))
	for i := 0; i < N; i++ {
		sb.WriteString(fmt.Sprintf("%s\n", A[i]))
	}
	for i := 0; i < M; i++ {
		sb.WriteString(fmt.Sprintf("%s %d\n", B[i], P[i]))
	}
	ok, order := solveCase(N, M, K, A, B, P)
	var out string
	if !ok {
		out = "NO\n"
	} else {
		var b strings.Builder
		b.WriteString("YES\n")
		for i, v := range order {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%d", v))
		}
		b.WriteByte('\n')
		out = b.String()
	}
	return testCase{input: sb.String(), output: out}
}

func randPattern(rng *rand.Rand, K int, allowUnderscore bool) string {
	b := make([]byte, K)
	for i := 0; i < K; i++ {
		if allowUnderscore && rng.Intn(3) == 0 {
			b[i] = '_'
		} else {
			b[i] = byte('a' + rng.Intn(3))
		}
	}
	return string(b)
}

func randomCase(rng *rand.Rand) testCase {
	K := rng.Intn(2) + 1
	N := rng.Intn(3) + 1
	M := rng.Intn(3) + 1
	set := make(map[string]bool)
	A := make([]string, 0, N)
	for len(A) < N {
		p := randPattern(rng, K, true)
		if !set[p] {
			set[p] = true
			A = append(A, p)
		}
	}
	B := make([]string, M)
	P := make([]int, M)
	for i := 0; i < M; i++ {
		B[i] = randPattern(rng, K, false)
		P[i] = rng.Intn(N) + 1
	}
	return buildCase(N, M, K, A, B, P)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase(1, 1, 1, []string{"a"}, []string{"a"}, []int{1}))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type item struct {
	idx   int
	limit int64
}

type maxHeap []item

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].limit > h[j].limit }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func canReachMax(h []int64, k int) bool {
	n := len(h)
	k--
	maxH := h[0]
	for _, v := range h {
		if v > maxH {
			maxH = v
		}
	}
	if h[k] == maxH {
		return true
	}
	visited := make([]bool, n)
	hp := &maxHeap{}
	heap.Init(hp)
	heap.Push(hp, item{idx: k, limit: h[k]})
	visited[k] = true
	for hp.Len() > 0 {
		cur := heap.Pop(hp).(item)
		if h[cur.idx] == maxH {
			return true
		}
		for _, nx := range []int{cur.idx - 1, cur.idx + 1} {
			if nx < 0 || nx >= n || visited[nx] {
				continue
			}
			newLimit := cur.limit - abs64(h[cur.idx]-h[nx])
			if newLimit >= h[nx] {
				visited[nx] = true
				heap.Push(hp, item{idx: nx, limit: newLimit})
			}
		}
	}
	return false
}

type testCase struct {
	n int
	k int
	h []int64
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(30) + 1
	k := rng.Intn(n) + 1
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		// heights between 1 and 50
		h[i] = int64(rng.Intn(50) + 1)
	}
	return testCase{n: n, k: k, h: h}
}

func buildInput(cases []testCase) (string, []string) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	exp := make([]string, len(cases))
	for i, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for j, v := range tc.h {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		if canReachMax(tc.h, tc.k) {
			exp[i] = "YES"
		} else {
			exp[i] = "NO"
		}
	}
	return sb.String(), exp
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/2126C_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		got := strings.ToUpper(lines[i])
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\nn=%d k=%d h=%v\n", i+1, exp, got, cases[i].n, cases[i].k, cases[i].h)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

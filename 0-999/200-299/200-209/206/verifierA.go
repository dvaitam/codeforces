package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type seq struct {
	value int64
	id    int
}

type seqHeap []seq

func (h seqHeap) Len() int            { return len(h) }
func (h seqHeap) Less(i, j int) bool  { return h[i].value < h[j].value }
func (h seqHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *seqHeap) Push(x interface{}) { *h = append(*h, x.(seq)) }
func (h *seqHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	k := make([]int, n)
	xi := make([]int64, n)
	yi := make([]int64, n)
	mi := make([]int64, n)
	curr := make([]int64, n)
	pos := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		var a1 int64
		fmt.Fscan(in, &k[i], &a1, &xi[i], &yi[i], &mi[i])
		curr[i] = a1
		pos[i] = 1
		total += k[i]
	}
	h := &seqHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		if k[i] > 0 {
			heap.Push(h, seq{value: curr[i], id: i})
		}
	}
	var badCount int64
	var prev int64
	first := true
	for h.Len() > 0 {
		s := heap.Pop(h).(seq)
		v := s.value
		i := s.id
		if first {
			prev = v
			first = false
		} else {
			if prev > v {
				badCount++
			}
			prev = v
		}
		pos[i]++
		if pos[i] < k[i] {
			next := (curr[i]*xi[i] + yi[i]) % mi[i]
			curr[i] = next
			heap.Push(h, seq{value: next, id: i})
		}
	}
	return fmt.Sprintf("%d\n", badCount)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	lines := make([]string, n+1)
	lines[0] = fmt.Sprintf("%d", n)
	for i := 0; i < n; i++ {
		k := rng.Intn(4) + 1
		m := rng.Int63n(1000) + 1
		a1 := rng.Int63n(m)
		xi := rng.Int63n(1000) + 1
		yi := rng.Int63n(1000)
		lines[i+1] = fmt.Sprintf("%d %d %d %d %d", k, a1, xi, yi, m)
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		expect := strings.TrimSpace(solveA(tc))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLine := strings.TrimSpace(strings.SplitN(got, "\n", 2)[0])
		if gotLine != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, gotLine, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveE(n, m, d int, a []int) string {
	h := &minHeap{}
	heap.Init(h)
	sumPrev := 0
	best := 0
	for i := 0; i < n; i++ {
		cand := a[i] + sumPrev - d*(i+1)
		if cand > best {
			best = cand
		}
		if a[i] > 0 {
			heap.Push(h, a[i])
			sumPrev += a[i]
			if h.Len() > m-1 {
				sumPrev -= heap.Pop(h).(int)
			}
		}
	}
	return fmt.Sprint(best)
}

func genCases() []string {
	rand.Seed(5)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		m := rand.Intn(n) + 1
		d := rand.Intn(5) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(21) - 10
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, d))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(a[j]))
		}
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n, m, d int
		fmt.Sscan(lines[1], &n, &m, &d)
		parts := strings.Fields(lines[2])
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &a[j])
		}
		want := solveE(n, m, d, a)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

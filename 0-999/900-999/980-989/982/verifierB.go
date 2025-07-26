package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

type row struct {
	w   int
	idx int
}

type maxHeap []row

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].w > h[j].w }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(row)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solve(n int, w []int, s string) string {
	pairs := make([]row, n)
	for i := 0; i < n; i++ {
		pairs[i] = row{w[i], i + 1}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].w < pairs[j].w })
	h := &maxHeap{}
	heap.Init(h)
	p := 0
	res := make([]int, 0, 2*n)
	for _, ch := range s {
		if ch == '0' {
			r := pairs[p]
			p++
			heap.Push(h, r)
			res = append(res, r.idx)
		} else {
			r := heap.Pop(h).(row)
			res = append(res, r.idx)
		}
	}
	var b strings.Builder
	for i, v := range res {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	return b.String()
}

func generateTests() []testCase {
	rand.Seed(2)
	var tests []testCase
	// simple edge cases
	tests = append(tests, testCase{
		input:  "1\n1\n01\n",
		output: "1 1",
	})
	// random cases
	for len(tests) < 120 {
		n := rand.Intn(10) + 1
		weights := rand.Perm(2*n + 5)
		w := make([]int, n)
		for i := 0; i < n; i++ {
			w[i] = weights[i] + 1
		}
		sBytes := make([]byte, 2*n)
		for i := 0; i < n; i++ {
			sBytes[i] = '0'
			sBytes[i+n] = '1'
		}
		rand.Shuffle(2*n, func(i, j int) { sBytes[i], sBytes[j] = sBytes[j], sBytes[i] })
		s := string(sBytes)
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for i, val := range w {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteString("\n")
		b.WriteString(s)
		b.WriteString("\n")
		tests = append(tests, testCase{
			input:  b.String(),
			output: solve(n, w, s),
		})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

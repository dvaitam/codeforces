package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type minHeap []int64

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func expectedOutput(n, k int, sVals, mVals []int64) string {
	res := make([]int64, n)
	h := &minHeap{}
	heap.Init(h)
	free := k
	type job struct {
		idx int
		m   int64
	}
	queue := make([]job, 0)
	process := func(t int64) {
		for h.Len() > 0 {
			ft := (*h)[0]
			if ft > t {
				break
			}
			heap.Pop(h)
			if len(queue) > 0 {
				j := queue[0]
				queue = queue[1:]
				start := ft
				finish := start + j.m
				res[j.idx] = finish
				heap.Push(h, finish)
			} else {
				free++
			}
		}
	}
	for i := 0; i < n; i++ {
		ti := sVals[i]
		process(ti)
		if free > 0 {
			free--
			start := ti
			finish := start + mVals[i]
			res[i] = finish
			heap.Push(h, finish)
		} else {
			queue = append(queue, job{i, mVals[i]})
		}
	}
	for len(queue) > 0 {
		ft := heap.Pop(h).(int64)
		j := queue[0]
		queue = queue[1:]
		start := ft
		finish := start + j.m
		res[j.idx] = finish
		heap.Push(h, finish)
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d", res[i])
		if i+1 < n {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var tcases int
	if _, err := fmt.Fscan(reader, &tcases); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= tcases; caseNum++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		sVals := make([]int64, n)
		mVals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &sVals[i], &mVals[i])
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d %d\n", sVals[i], mVals[i])
		}
		want := expectedOutput(n, k, sVals, mVals)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tcases)
}

package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type basket struct {
	cnt   int
	dist2 int
	idx   int
}

type basketHeap []basket

func (h basketHeap) Len() int { return len(h) }
func (h basketHeap) Less(i, j int) bool {
	if h[i].cnt != h[j].cnt {
		return h[i].cnt < h[j].cnt
	}
	if h[i].dist2 != h[j].dist2 {
		return h[i].dist2 < h[j].dist2
	}
	return h[i].idx < h[j].idx
}
func (h basketHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *basketHeap) Push(x interface{}) { *h = append(*h, x.(basket)) }
func (h *basketHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(n, m int) []int {
	h := make(basketHeap, 0, m)
	center := m + 1
	for i := 1; i <= m; i++ {
		d2 := center - 2*i
		if d2 < 0 {
			d2 = -d2
		}
		h = append(h, basket{cnt: 0, dist2: d2, idx: i})
	}
	heap.Init(&h)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		b := heap.Pop(&h).(basket)
		res[i] = b.idx
		b.cnt++
		heap.Push(&h, b)
	}
	return res
}

func solveBFromInput(input string) ([]int, int, error) {
	r := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return nil, 0, err
	}
	return solveB(n, m), n, nil
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	m := rng.Intn(50) + 1
	return fmt.Sprintf("%d %d\n", n, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect, n, err := solveBFromInput(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad generated test: %v\n", err)
			os.Exit(1)
		}
		gotStr, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		fields := strings.Fields(gotStr)
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, n, len(fields), tc)
			os.Exit(1)
		}
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil || v != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at ball %d: expected %d got %s\ninput:\n%s", i+1, j+1, expect[j], f, tc)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

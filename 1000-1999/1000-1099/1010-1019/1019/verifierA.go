package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Pair struct{ first, second int }

type LaterHeap struct {
	data  [][]int64
	items []Pair
}

func (h LaterHeap) Len() int { return len(h.items) }
func (h LaterHeap) Less(i, j int) bool {
	a := h.items[i]
	b := h.items[j]
	return h.data[a.first][a.second] < h.data[b.first][b.second]
}
func (h LaterHeap) Swap(i, j int)       { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *LaterHeap) Push(x interface{}) { h.items = append(h.items, x.(Pair)) }
func (h *LaterHeap) Pop() interface{} {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[:n-1]
	return x
}

func compute(n, m int, party []int, cost []int64) int64 {
	cp := make([][]int64, m)
	for i := 0; i < n; i++ {
		cp[party[i]] = append(cp[party[i]], cost[i])
	}
	ps := make([][]int64, m)
	for i := 0; i < m; i++ {
		sort.Slice(cp[i], func(a, b int) bool { return cp[i][a] < cp[i][b] })
		ps[i] = make([]int64, len(cp[i]))
		for j, v := range cp[i] {
			if j == 0 {
				ps[i][j] = v
			} else {
				ps[i][j] = ps[i][j-1] + v
			}
		}
	}
	ans := int64(1 << 62)
	for b := 0; b <= n; b++ {
		tot := len(cp[0])
		cst := int64(0)
		h := &LaterHeap{data: cp}
		heap.Init(h)
		for i := 1; i < m; i++ {
			sz := len(cp[i])
			need := sz - b
			if need < 0 {
				need = 0
			}
			if need > 0 {
				cst += ps[i][need-1]
				tot += need
			}
			if need < sz {
				heap.Push(h, Pair{i, need})
			}
		}
		for tot <= b && h.Len() > 0 {
			top := heap.Pop(h).(Pair)
			cst += cp[top.first][top.second]
			if top.second+1 < len(cp[top.first]) {
				heap.Push(h, Pair{top.first, top.second + 1})
			}
			tot++
		}
		if tot > b && cst < ans {
			ans = cst
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(4) + 2
		party := make([]int, n)
		cost := make([]int64, n)
		for j := 0; j < n; j++ {
			party[j] = rng.Intn(m)
			cost[j] = int64(rng.Intn(20) + 1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", party[j]+1, cost[j])
		}
		input := sb.String()
		expected := compute(n, m, party, cost)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		val, err2 := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err2 != nil || val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

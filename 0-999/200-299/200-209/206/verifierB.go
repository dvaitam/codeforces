package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// -------- reference solution --------
type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	N := 2 * n
	pairs := make([]struct{ l, j int }, N)
	for j := 1; j <= N; j++ {
		aj := a[(j-1)%n]
		l := j - aj
		if l < 1 {
			l = 1
		}
		pairs[j-1] = struct{ l, j int }{l, j}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].l < pairs[j].l })
	R := make([]int, N+2)
	h := &maxHeap{}
	heap.Init(h)
	idx := 0
	for i := 1; i <= N; i++ {
		for idx < N && pairs[idx].l <= i {
			heap.Push(h, pairs[idx].j)
			idx++
		}
		for h.Len() > 0 {
			top := (*h)[0]
			if top <= i {
				heap.Pop(h)
			} else {
				break
			}
		}
		if h.Len() == 0 {
			R[i] = i
		} else {
			R[i] = (*h)[0]
		}
	}
	R[N+1] = N + 1
	maxLg := 0
	for (1 << maxLg) <= N {
		maxLg++
	}
	up := make([][]int, maxLg)
	up[0] = make([]int, N+2)
	for i := 1; i <= N; i++ {
		up[0][i] = R[i]
	}
	for k := 1; k < maxLg; k++ {
		up[k] = make([]int, N+2)
		for i := 1; i <= N; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
	var total int64
	for s := 1; s <= n; s++ {
		target := s + n - 1
		cur := s
		var cnt int64
		for k := maxLg - 1; k >= 0; k-- {
			nxt := up[k][cur]
			if nxt < target {
				cur = nxt
				cnt += 1 << k
			}
		}
		if cur < target {
			cnt++
		}
		total += cnt
	}
	return fmt.Sprintf("%d\n", total)
}

// -------- runner --------
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
	n := rng.Intn(5) + 1
	nums := make([]string, n+1)
	nums[0] = fmt.Sprintf("%d", n)
	for i := 0; i < n; i++ {
		nums[i+1] = fmt.Sprintf("%d", rng.Intn(10)+1)
	}
	return strings.Join(nums, " \n") + "\n"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		expect := strings.TrimSpace(solveB(tc))
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

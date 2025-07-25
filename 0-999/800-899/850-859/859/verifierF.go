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

// Item and heaps from 859F solution

type Item struct {
	t   int64
	f   int64
	idx int
}

type minHeap []Item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].t < h[j].t }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *minHeap) Pop() any          { n := len(*h); v := (*h)[n-1]; *h = (*h)[:n-1]; return v }

type maxHeap []Item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].f > h[j].f }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *maxHeap) Pop() any          { n := len(*h); v := (*h)[n-1]; *h = (*h)[:n-1]; return v }

func minimalShirts(n int, C int64, single []int64, multi []int64) int64 {
	multiExt := make([]int64, n+1)
	for i := 1; i < n; i++ {
		multiExt[i] = multi[i-1]
	}
	S1 := make([]int64, n+1)
	S2 := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		S1[i] = S1[i-1] + single[i-1]
	}
	for i := 1; i <= n; i++ {
		S2[i] = S2[i-1] + multiExt[i]
	}
	PS := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		PS[i] = S1[i] + S2[i-1]
	}
	dp := make([]int64, n+1)
	satMax := int64(-1 << 62)
	uheap := &maxHeap{}
	theap := &minHeap{}
	heap.Init(uheap)
	heap.Init(theap)
	active := make([]bool, n+1)

	for i := 1; i <= n; i++ {
		F1 := dp[i-1] - PS[i-1] - multiExt[i-1]
		T := PS[i-1] + multiExt[i-1] + C
		it := Item{t: T, f: F1, idx: i}
		heap.Push(uheap, it)
		heap.Push(theap, it)
		active[i] = true
		for theap.Len() > 0 && (*theap)[0].t <= PS[i] {
			v := heap.Pop(theap).(Item)
			if !active[v.idx] {
				continue
			}
			active[v.idx] = false
			val := dp[v.idx-1] + C
			if val > satMax {
				satMax = val
			}
		}
		for uheap.Len() > 0 && !active[(*uheap)[0].idx] {
			heap.Pop(uheap)
		}
		bestUnsat := int64(-1 << 62)
		if uheap.Len() > 0 {
			bestUnsat = (*uheap)[0].f + PS[i]
		}
		if satMax > bestUnsat {
			dp[i] = satMax
		} else {
			dp[i] = bestUnsat
		}
	}
	return dp[n]
}

func solveF(n int, C int64, values []int64) int64 {
	single := make([]int64, n)
	multi := make([]int64, n-1)
	for i := 0; i < n; i++ {
		single[i] = values[2*i]
		if i < n-1 {
			multi[i] = values[2*i+1]
		}
	}
	return minimalShirts(n, C, single, multi)
}

type CaseF struct {
	input    string
	expected int64
}

func generateCaseF(rng *rand.Rand) CaseF {
	n := rng.Intn(20) + 1
	C := rng.Int63n(1_000_000) + 1
	values := make([]int64, 2*n-1)
	for i := range values {
		values[i] = rng.Int63n(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, C))
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return CaseF{sb.String(), solveF(n, C, values)}
}

func runCase(exe, input string, expected int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	values := []int64{1}
	cases := []CaseF{{"1 1\n1\n", solveF(1, 1, values)}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseF(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

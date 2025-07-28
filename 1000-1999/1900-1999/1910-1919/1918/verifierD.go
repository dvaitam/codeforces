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

type Pair struct {
	val int64
	idx int
}

type MinHeap []Pair

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(Pair)) }
func (h *MinHeap) Pop() any {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

func minBlockedSum(a []int64, M int64) int64 {
	n := len(a)
	prefix := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}
	prefix[n+1] = prefix[n]

	arr := append(a, 0)

	dp := make([]int64, n+2)
	h := &MinHeap{}
	heap.Init(h)
	heap.Push(h, Pair{0, 0})
	L := 0
	for i := 1; i <= n+1; i++ {
		thresh := prefix[i-1] - M
		for L <= i-1 && prefix[L] < thresh {
			L++
		}
		for h.Len() > 0 && (*h)[0].idx < L {
			heap.Pop(h)
		}
		minVal := int64(1 << 62)
		if h.Len() > 0 {
			minVal = (*h)[0].val
		}
		dp[i] = arr[i-1] + minVal
		heap.Push(h, Pair{dp[i], i})
	}
	return dp[n+1]
}

func solveArr(a []int64) int64 {
	var mx, sum int64
	for _, v := range a {
		if v > mx {
			mx = v
		}
		sum += v
	}
	l, r := mx, sum
	for l < r {
		mid := (l + r) / 2
		if minBlockedSum(a, mid) <= mid {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func runCase(bin string, a []int64) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", solveArr(a))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := make([][]int64, 0, 100)
	cases = append(cases, []int64{1})
	for i := 0; i < 99; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(1000) + 1
		}
		cases = append(cases, arr)
	}
	for i, a := range cases {
		if err := runCase(bin, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

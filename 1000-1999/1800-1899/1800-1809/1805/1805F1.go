package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/big"
	"os"
	"sort"
)

type Pair struct {
	sum *big.Int
	i   int
	j   int
}

type PairHeap []*Pair

func (h PairHeap) Len() int            { return len(h) }
func (h PairHeap) Less(a, b int) bool  { return h[a].sum.Cmp(h[b].sum) < 0 }
func (h PairHeap) Swap(a, b int)       { h[a], h[b] = h[b], h[a] }
func (h *PairHeap) Push(x interface{}) { *h = append(*h, x.(*Pair)) }
func (h *PairHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func F(arr []*big.Int) []*big.Int {
	n := len(arr)
	sort.Slice(arr, func(i, j int) bool { return arr[i].Cmp(arr[j]) < 0 })
	h := &PairHeap{}
	heap.Init(h)
	for i := 0; i < n-1; i++ {
		sum := new(big.Int).Add(arr[i], arr[i+1])
		heap.Push(h, &Pair{sum: sum, i: i, j: i + 1})
	}
	res := make([]*big.Int, n-1)
	for k := 0; k < n-1; k++ {
		p := heap.Pop(h).(*Pair)
		res[k] = p.sum
		if p.j+1 < n {
			sum := new(big.Int).Add(arr[p.i], arr[p.j+1])
			heap.Push(h, &Pair{sum: sum, i: p.i, j: p.j + 1})
		}
	}
	return res
}

const MOD int64 = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		arr[i] = big.NewInt(v)
	}
	for len(arr) > 1 {
		arr = F(arr)
	}
	mod := big.NewInt(MOD)
	ans := new(big.Int).Mod(arr[0], mod)
	fmt.Println(ans.Int64())
}

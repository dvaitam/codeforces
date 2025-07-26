package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Int64Heap []int64

func (h Int64Heap) Len() int           { return len(h) }
func (h Int64Heap) Less(i, j int) bool { return h[i] < h[j] }
func (h Int64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *Int64Heap) Push(x any)        { *h = append(*h, x.(int64)) }
func (h *Int64Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var m int64
		var k int
		fmt.Fscan(reader, &n, &m, &k)
		h := &Int64Heap{}
		heap.Init(h)
		var sum int64
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(reader, &v)
			heap.Push(h, v)
			sum += v
		}
		for i := 0; i < k; i++ {
			x := heap.Pop(h).(int64)
			y := heap.Pop(h).(int64)
			g := gcd(x, y)
			sum -= x + y - g
			heap.Push(h, g)
		}
		fmt.Fprintln(writer, sum)
	}
}

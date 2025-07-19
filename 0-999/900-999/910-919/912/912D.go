package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Node represents an element in the max-heap
type Node struct {
   now     int64
   x, y    int64
   shu     int64
   flag    bool
}

// MaxHeap implements a max-heap of Nodes based on now
type MaxHeap []Node

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].now > h[j].now }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Node)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

var n, m, r, k int64

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

// get computes the number of r-by-r submatrices containing cell at (x,y)
func get(x, y int64) int64 {
   lx := max(1, x-r+1)
   ly := max(1, y-r+1)
   x = min(x, n-r+1)
   y = min(y, m-r+1)
   return (x-lx+1) * (y-ly+1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m, &r, &k)

   number := (n - r + 1) * (m - r + 1)
   MAX := get((n+1)/2, (m+1)/2)
   var a, b int64
   if n >= 2*r {
       a = n - 2*r + 2
   } else {
       a = 2*r - n
   }
   if m >= 2*r {
       b = m - 2*r + 2
   } else {
       b = 2*r - m
   }
   hen := min(n-r+1, r)
   shu := min(m-r+1, r)

   h := &MaxHeap{}
   heap.Init(h)
   heap.Push(h, Node{now: MAX, x: a, y: b, shu: shu, flag: true})

   var ans int64
   for {
       p := heap.Pop(h).(Node)
       cnt := p.x * p.y
       if cnt >= k {
           ans += p.now * k
           break
       }
       k -= cnt
       ans += p.now * cnt
       if p.flag && p.now-hen >= 1 {
           heap.Push(h, Node{now: p.now - hen, x: p.x, y: 2, shu: p.shu - 1, flag: true})
       }
       if p.now-p.shu >= 1 {
           heap.Push(h, Node{now: p.now - p.shu, x: 2, y: p.y, shu: p.shu, flag: false})
       }
   }

   fmt.Fprintf(writer, "%.15f", float64(ans)/float64(number))
}

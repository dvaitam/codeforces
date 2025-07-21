package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// DSU for union-find
type DSU struct {
   p, sz []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   sz := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
       sz[i] = 1
   }
   return &DSU{p: p, sz: sz}
}

func (d *DSU) Find(x int) int {
   for d.p[x] != x {
       d.p[x] = d.p[d.p[x]]
       x = d.p[x]
   }
   return x
}

func (d *DSU) Union(a, b int) {
   ra := d.Find(a)
   rb := d.Find(b)
   if ra == rb {
       return
   }
   if d.sz[ra] < d.sz[rb] {
       ra, rb = rb, ra
   }
   d.p[rb] = ra
   d.sz[ra] += d.sz[rb]
}

// min-heap of ints
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   dsu := NewDSU(n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       dsu.Union(u, v)
   }
   // count component sizes
   compCap := make(map[int]int)
   for i := 1; i <= n; i++ {
       r := dsu.Find(i)
       compCap[r]++
   }
   // build heap of caps and sum
   h := &IntHeap{}
   heap.Init(h)
   sumCaps := int64(0)
   D := 0
   for _, sz := range compCap {
       cap := sz
       if cap > k {
           cap = k
       }
       heap.Push(h, cap)
       sumCaps += int64(cap)
       D++
   }
   if D <= 1 {
       fmt.Println(0)
       return
   }
   // merge until feasible: sumCaps >= 2*(D-1)
   roads := 0
   for sumCaps < int64(2*(D-1)) {
       // take two smallest
       c1 := heap.Pop(h).(int)
       c2 := heap.Pop(h).(int)
       newC := c1 + c2
       loss := 0
       if newC > k {
           loss = newC - k
           newC = k
       }
       sumCaps -= int64(loss)
       D--
       roads++
       heap.Push(h, newC)
   }
   fmt.Println(roads)
}

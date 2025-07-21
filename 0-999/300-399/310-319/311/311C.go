package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
)

const INF64 = 9e18

// Dijkstra node for residues
type dNode struct {
   d int64
   r int
}
// min-heap for dNode
type dHeap []dNode
func (h dHeap) Len() int            { return len(h) }
func (h dHeap) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h dHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *dHeap) Push(x interface{}) { *h = append(*h, x.(dNode)) }
func (h *dHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// Treasure heap entry
type tEntry struct {
   valNeg  int64
   idx     int
   version int
}
// max-heap via valNeg (smaller is better), tie by smaller idx
type tHeap []tEntry
func (h tHeap) Len() int { return len(h) }
func (h tHeap) Less(i, j int) bool {
   if h[i].valNeg != h[j].valNeg {
       return h[i].valNeg < h[j].valNeg
   }
   return h[i].idx < h[j].idx
}
func (h tHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *tHeap) Push(x interface{}) { *h = append(*h, x.(tEntry)) }
func (h *tHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var hVal int64
   var n, m int
   var k0 int64
   if _, err := fmt.Fscan(in, &hVal, &n, &m, &k0); err != nil {
       if err == io.EOF {
           return
       }
       panic(err)
   }
   pos := make([]int64, n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pos[i], &c[i])
   }
   methods := []int64{k0}
   exists := make([]bool, n)
   reachable := make([]bool, n)
   version := make([]int, n)
   for i := 0; i < n; i++ {
       exists[i] = true
       reachable[i] = false
       version[i] = 0
   }
   // treasure heap
   th := &tHeap{}
   heap.Init(th)

   // function to recompute reachability
   var recompute func()
   recompute = func() {
       // find minimal method
       a0 := methods[0]
       for _, v := range methods {
           if v < a0 {
               a0 = v
           }
       }
       a0i := int(a0)
       // Dijkstra on residues mod a0
       dist := make([]int64, a0i)
       for i := 0; i < a0i; i++ {
           dist[i] = INF64
       }
       dist[0] = 0
       dh := &dHeap{{d: 0, r: 0}}
       heap.Init(dh)
       for dh.Len() > 0 {
           nd := heap.Pop(dh).(dNode)
           dcur, r := nd.d, nd.r
           if dcur != dist[r] {
               continue
           }
           for _, w := range methods {
               nr := (r + int(w%a0)) % a0i
               ndist := dcur + w
               if ndist < dist[nr] {
                   dist[nr] = ndist
                   heap.Push(dh, dNode{d: ndist, r: nr})
               }
           }
       }
       // scan treasures
       for i := 0; i < n; i++ {
           if !exists[i] || reachable[i] {
               continue
           }
           delta := pos[i] - 1
           r := int(delta % a0)
           if dist[r] <= delta {
               reachable[i] = true
               heap.Push(th, tEntry{valNeg: -c[i], idx: i, version: version[i]})
           }
       }
   }
   // initial
   recompute()
   // process operations
   for opi := 0; opi < m; opi++ {
       var typ int
       fmt.Fscan(in, &typ)
       if typ == 1 {
           var x int64
           fmt.Fscan(in, &x)
           methods = append(methods, x)
           recompute()
       } else if typ == 2 {
           var xi int
           var y int64
           fmt.Fscan(in, &xi, &y)
           xi--
           if !exists[xi] {
               continue
           }
           c[xi] -= y
           version[xi]++
           if reachable[xi] {
               heap.Push(th, tEntry{valNeg: -c[xi], idx: xi, version: version[xi]})
           }
       } else if typ == 3 {
           // pop best
           var ans int64
           for th.Len() > 0 {
               e := heap.Pop(th).(tEntry)
               i := e.idx
               if exists[i] && reachable[i] && version[i] == e.version {
                   ans = -e.valNeg
                   exists[i] = false
                   break
               }
           }
           if ans == 0 {
               fmt.Fprintln(out, 0)
           } else {
               fmt.Fprintln(out, ans)
           }
       }
   }
}

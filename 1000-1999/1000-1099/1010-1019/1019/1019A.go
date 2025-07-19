package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
   "sort"
)

// Pair holds indices into cp and ps
type Pair struct {
   first, second int
}

// LaterHeap is a min-heap of Pair based on cp values
type LaterHeap struct {
   data [][]int64
   items []Pair
}

func (h LaterHeap) Len() int { return len(h.items) }
func (h LaterHeap) Less(i, j int) bool {
   a := h.items[i]
   b := h.items[j]
   return h.data[a.first][a.second] < h.data[b.first][b.second]
}
func (h LaterHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *LaterHeap) Push(x interface{}) {
   h.items = append(h.items, x.(Pair))
}
func (h *LaterHeap) Pop() interface{} {
   old := h.items
   n := len(old)
   x := old[n-1]
   h.items = old[0 : n-1]
   return x
}

// Solver encapsulates the problem state
type Solver struct {
   n, m int
   cp   [][]int64
   ps   [][]int64
   ans  int64
}

// cost computes minimal cost to achieve more than bound votes
func (s *Solver) cost(bound int) int64 {
   var ret int64
   tota := len(s.cp[0])
   // initialize heap
   h := &LaterHeap{data: s.cp}
   heap.Init(h)
   // for each other candidate
   for i := 1; i < s.m; i++ {
       sz := len(s.cp[i])
       here := sz - bound
       if here < 0 {
           here = 0
       }
       if here > 0 {
           ret += s.ps[i][here-1]
           tota += here
       }
       if here < sz {
           heap.Push(h, Pair{i, here})
       }
   }
   // bribe cheapest until we have more than bound
   for tota <= bound && h.Len() > 0 {
       top := heap.Pop(h).(Pair)
       ret += s.cp[top.first][top.second]
       if top.second+1 < len(s.cp[top.first]) {
           heap.Push(h, Pair{top.first, top.second + 1})
       }
       tota++
   }
   if tota > bound {
       return ret
   }
   // cannot reach
   return 1<<62
}

func (s *Solver) Solve() {
   // sort and prefix sums
   for i := 0; i < s.m; i++ {
       sort.Slice(s.cp[i], func(a, b int) bool { return s.cp[i][a] < s.cp[i][b] })
       s.ps[i] = make([]int64, len(s.cp[i]))
       for j, v := range s.cp[i] {
           if j == 0 {
               s.ps[i][j] = v
           } else {
               s.ps[i][j] = s.ps[i][j-1] + v
           }
       }
   }
   s.ans = 1<<62
   // try all possible bounds
   for b := 0; b <= s.n; b++ {
       cur := s.cost(b)
       if cur > s.ans {
           break
       }
       if cur < s.ans {
           s.ans = cur
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil && err != io.EOF {
       return
   }
   s := &Solver{n: n, m: m}
   s.cp = make([][]int64, m)
   s.ps = make([][]int64, m)
   for i := 0; i < n; i++ {
       var p int
       var c int64
       fmt.Fscan(reader, &p, &c)
       p--
       s.cp[p] = append(s.cp[p], c)
   }
   s.Solve()
   // output answer
   fmt.Println(s.ans)
}

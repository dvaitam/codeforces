package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// IntHeap is a min-heap of ints.
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
// Peek returns the minimum element without removing it
func (h *IntHeap) Peek() int {
   return (*h)[0]
}

// State holds the multiset of depth->count and total flows
type State struct {
   m     map[int]int64
   heap  *IntHeap
   total int64
}

// merge merges other into s (other may be discarded)
func (s *State) merge(other *State) {
   if other == nil || other.m == nil {
       return
   }
   if len(s.m) < len(other.m) {
       // swap
       s.m, other.m = other.m, s.m
       s.heap, other.heap = other.heap, s.heap
       s.total, other.total = other.total, s.total
   }
   // merge other into s
   for d, cnt := range other.m {
       if cnt == 0 {
           continue
       }
       if _, ok := s.m[d]; ok {
           s.m[d] += cnt
       } else {
           s.m[d] = cnt
           heap.Push(s.heap, d)
       }
       s.total += cnt
   }
   // clear other to release memory
   other.m = nil
   other.heap = nil
   other.total = 0
}

var (
   n, sRoot int
   k     []int64
   adj   [][]int
   depth []int
)

// dfs returns state for subtree rooted at u
func dfs(u, p int) *State {
   st := &State{m: make(map[int]int64), heap: &IntHeap{}, total: 0}
   heap.Init(st.heap)
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       depth[v] = depth[u] + 1
       child := dfs(v, u)
       st.merge(child)
   }
   // add own endpoint (except root)
   if u != sRoot {
       d := depth[u]
       cnt := k[u]
       if cnt > 0 {
           if _, ok := st.m[d]; ok {
               st.m[d] += cnt
           } else {
               st.m[d] = cnt
               heap.Push(st.heap, d)
           }
           st.total += cnt
       }
   }
   // enforce capacity k[u]
   cap := k[u]
   for st.total > cap {
       // get smallest depth
       d := st.heap.Peek()
       cnt := st.m[d]
       toRemove := st.total - cap
       if cnt <= toRemove {
           // remove all
           heap.Pop(st.heap)
           delete(st.m, d)
           st.total -= cnt
       } else {
           // partial remove
           st.m[d] = cnt - toRemove
           st.total -= toRemove
           break
       }
   }
   return st
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   k = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &k[i])
   }
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u, v := 0, 0
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   fmt.Fscan(reader, &sRoot)
   depth = make([]int, n+1)
   depth[sRoot] = 0
   st := dfs(sRoot, 0)
   var ans int64 = 0
   for d, cnt := range st.m {
       ans += int64(d) * cnt * 2
   }
   fmt.Fprintln(writer, ans)
}

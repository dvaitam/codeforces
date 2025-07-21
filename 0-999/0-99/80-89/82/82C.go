package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Div represents a division with priority ai, arrival time t at a node, and its id.
type Div struct {
   ai int
   t  int
   id int
}

// DivHeap is a min-heap of Div by ai.
type DivHeap []Div
func (h DivHeap) Len() int            { return len(h) }
func (h DivHeap) Less(i, j int) bool  { return h[i].ai < h[j].ai }
func (h DivHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *DivHeap) Push(x interface{}) { *h = append(*h, x.(Div)) }
func (h *DivHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// schedule performs batch scheduling on an edge with capacity cap.
// divs: divisions arriving at this node, with release time divs[i].t.
// Returns divisions with updated t = arrival time at parent (departure day + 1).
func schedule(divs []Div, cap int) []Div {
   m := len(divs)
   // sort by release time
   sort.Slice(divs, func(i, j int) bool { return divs[i].t < divs[j].t })
   var pq DivHeap
   heap.Init(&pq)
   result := make([]Div, 0, m)
   idx := 0
   day := 0
   for len(result) < m {
       // advance day if no ready tasks
       if pq.Len() == 0 && idx < m && divs[idx].t > day {
           day = divs[idx].t
       }
       // add released tasks
       for idx < m && divs[idx].t <= day {
           heap.Push(&pq, divs[idx])
           idx++
       }
       // depart up to cap tasks
       for k := 0; k < cap && pq.Len() > 0; k++ {
           d := heap.Pop(&pq).(Div)
           d.t = day + 1
           result = append(result, d)
       }
       day++
   }
   return result
}

var (
   n            int
   a            []int
   children     [][]int
   capToParent  []int
   answer       []int
)

// dfs returns divisions arriving at node u (before edge to parent), with t at arrival at u.
func dfs(u int) []Div {
   // start with own division at u, arrival t=0
   divs := []Div{{ai: a[u], t: 0, id: u}}
   for _, v := range children[u] {
       // get divisions arrival at v
       sub := dfs(v)
       // schedule on edge v->u, capacity capToParent[v]
       scheduled := schedule(sub, capToParent[v])
       // these arrive at u
       divs = append(divs, scheduled...)
   }
   return divs
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build tree
   adj := make([][]struct{to, cap int}, n+1)
   for i := 0; i < n-1; i++ {
       var u, v, c int
       fmt.Fscan(reader, &u, &v, &c)
       adj[u] = append(adj[u], struct{to, cap int}{v, c})
       adj[v] = append(adj[v], struct{to, cap int}{u, c})
   }
   // prepare parent/children
   children = make([][]int, n+1)
   capToParent = make([]int, n+1)
   parent := make([]int, n+1)
   // BFS
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, e := range adj[u] {
           v, c := e.to, e.cap
           if v == parent[u] {
               continue
           }
           parent[v] = u
           capToParent[v] = c
           children[u] = append(children[u], v)
           queue = append(queue, v)
       }
   }
   answer = make([]int, n+1)
   // process
   rootDivs := dfs(1)
   for _, d := range rootDivs {
       answer[d.id] = d.t
   }
   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, answer[i])
   }
   fmt.Fprintln(writer)
}

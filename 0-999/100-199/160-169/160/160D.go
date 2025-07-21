package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   u, v, w, id int
}

// DSU for union-find
type DSU struct {
   p, r []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   r := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(a, b int) {
   a = d.Find(a)
   b = d.Find(b)
   if a == b {
       return
   }
   if d.r[a] < d.r[b] {
       a, b = b, a
   }
   d.p[b] = a
   if d.r[a] == d.r[b] {
       d.r[a]++
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
       edges[i].id = i
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })
   ans := make([]string, m)
   dsu := NewDSU(n)

   // process by weight
   for i := 0; i < m; {
       j := i
       for j < m && edges[j].w == edges[i].w {
           j++
       }
       // batch [i, j)
       // collect edges connecting different components
       type BE struct{u, v, id int}
       var be []BE
       for k := i; k < j; k++ {
           u := dsu.Find(edges[k].u)
           v := dsu.Find(edges[k].v)
           if u == v {
               ans[edges[k].id] = "none"
           } else {
               be = append(be, BE{u, v, edges[k].id})
           }
       }
       if len(be) > 0 {
           // map components to 0..sz-1
           comp := make(map[int]int)
           idx := 0
           for _, e := range be {
               if _, ok := comp[e.u]; !ok {
                   comp[e.u] = idx
                   idx++
               }
               if _, ok := comp[e.v]; !ok {
                   comp[e.v] = idx
                   idx++
               }
           }
           sz := idx
           // adjacency and count pairs
           type Adj struct{to, id int}
           adj := make([][]Adj, sz)
           pairCnt := make(map[int]int)
           key := func(a, b int) int {
               if a > b {
                   a, b = b, a
               }
               return a*sz + b
           }
           for _, e := range be {
               u := comp[e.u]
               v := comp[e.v]
               k := key(u, v)
               pairCnt[k]++
               adj[u] = append(adj[u], Adj{v, e.id})
               adj[v] = append(adj[v], Adj{u, e.id})
           }
           // bridge finding
           tin := make([]int, sz)
           low := make([]int, sz)
           vis := make([]bool, sz)
           timer := 0
           bridge := make(map[int]bool)
           var dfs func(u, pe int)
           dfs = func(u, pe int) {
               vis[u] = true
               timer++
               tin[u] = timer
               low[u] = timer
               for _, e := range adj[u] {
                   v := e.to
                   id := e.id
                   if id == pe {
                       continue
                   }
                   if vis[v] {
                       if tin[v] < low[u] {
                           low[u] = tin[v]
                       }
                   } else {
                       dfs(v, id)
                       if low[v] < low[u] {
                           low[u] = low[v]
                       }
                       // check bridge
                       if low[v] > tin[u] {
                           if pairCnt[key(u, v)] == 1 {
                               bridge[id] = true
                           }
                       }
                   }
               }
           }
           for u := 0; u < sz; u++ {
               if !vis[u] {
                   dfs(u, -1)
               }
           }
           // classify
           for _, e := range be {
               if bridge[e.id] {
                   ans[e.id] = "any"
               } else {
                   ans[e.id] = "at least one"
               }
           }
       }
       // unify
       for k := i; k < j; k++ {
           dsu.Union(edges[k].u, edges[k].v)
       }
       i = j
   }
   // output
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

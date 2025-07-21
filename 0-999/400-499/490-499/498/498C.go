package main

import (
   "bufio"
   "fmt"
   "os"
)

// Dinic max flow
type Edge struct {
   to, rev, cap int
}

type Dinic struct {
   n    int
   adj  [][]Edge
   lvl  []int
   it   []int
}

func NewDinic(n int) *Dinic {
   d := &Dinic{n: n, adj: make([][]Edge, n), lvl: make([]int, n), it: make([]int, n)}
   return d
}

func (d *Dinic) AddEdge(u, v, c int) {
   d.adj[u] = append(d.adj[u], Edge{to: v, rev: len(d.adj[v]), cap: c})
   d.adj[v] = append(d.adj[v], Edge{to: u, rev: len(d.adj[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := range d.lvl {
       d.lvl[i] = -1
   }
   queue := make([]int, 0, d.n)
   d.lvl[s] = 0
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, e := range d.adj[u] {
           if e.cap > 0 && d.lvl[e.to] < 0 {
               d.lvl[e.to] = d.lvl[u] + 1
               queue = append(queue, e.to)
               if e.to == t {
                   return true
               }
           }
       }
   }
   return d.lvl[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
   if u == t {
       return f
   }
   for i := d.it[u]; i < len(d.adj[u]); i++ {
       d.it[u] = i
       e := &d.adj[u][i]
       if e.cap > 0 && d.lvl[e.to] == d.lvl[u]+1 {
           ret := d.dfs(e.to, t, min(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.adj[e.to][e.rev].cap += ret
               return ret
           }
       }
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := range d.it {
           d.it[i] = 0
       }
       for {
           pushed := d.dfs(s, t, 1<<60)
           if pushed == 0 {
               break
           }
           flow += pushed
       }
   }
   return flow
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func factor(x int) map[int]int {
   m := make(map[int]int)
   d := 2
   for d*d <= x {
       for x%d == 0 {
           m[d]++
           x /= d
       }
       d++
   }
   if x > 1 {
       m[x]++
   }
   return m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   pairs := make([][2]int, m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       if u%2 == 0 {
           u, v = v, u
       }
       pairs[i][0], pairs[i][1] = u, v
   }
   cnt := make([]map[int]int, n+1)
   primes := make(map[int]struct{})
   for i := 1; i <= n; i++ {
       cnt[i] = factor(a[i])
       for p := range cnt[i] {
           primes[p] = struct{}{}
       }
   }
   total := 0
   const INF = 1000000000
   for p := range primes {
       // collect nodes
       oddIds := make(map[int]int)
       evenIds := make(map[int]int)
       oid, eid := 0, 0
       for i := 1; i <= n; i++ {
           if e, ok := cnt[i][p]; ok && e > 0 {
               if i%2 == 1 {
                   oid++
                   oddIds[i] = oid
               } else {
                   eid++
                   evenIds[i] = eid
               }
           }
       }
       if oid == 0 || eid == 0 {
           continue
       }
       N := 1 + oid + eid + 1
       src, sink := 0, N-1
       din := NewDinic(N)
       // source to odd
       for i, id := range oddIds {
           din.AddEdge(src, id, cnt[i][p])
       }
       // even to sink
       for i, id := range evenIds {
           din.AddEdge(oddIdsSize(oddIds)+id, sink, cnt[i][p])
       }
       // edges
       for _, pr := range pairs {
           u, v := pr[0], pr[1]
           ui, uok := oddIds[u]
           vi, vok := evenIds[v]
           if uok && vok {
               din.AddEdge(ui, oddIdsSize(oddIds)+vi, INF)
           }
       }
       flow := din.MaxFlow(src, sink)
       total += flow
   }
   fmt.Println(total)
}

func oddIdsSize(m map[int]int) int {
   // size of oddIds map
   return len(m)
}

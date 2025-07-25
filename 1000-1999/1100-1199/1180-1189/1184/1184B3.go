package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF64 = int64(1e18)

// Edge for Dinic
type Edge struct {
   to   int
   rev  int
   cap  int64
}

// Dinic max flow
type Dinic struct {
   N     int
   G     [][]Edge
   level []int
   iter  []int
}

func NewDinic(n int) *Dinic {
   g := make([][]Edge, n)
   lvl := make([]int, n)
   it := make([]int, n)
   return &Dinic{N: n, G: g, level: lvl, iter: it}
}

// Add edge u->v with capacity c
func (d *Dinic) AddEdge(u, v int, c int64) {
   d.G[u] = append(d.G[u], Edge{to: v, rev: len(d.G[v]), cap: c})
   d.G[v] = append(d.G[v], Edge{to: u, rev: len(d.G[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s int) {
   for i := range d.level {
       d.level[i] = -1
   }
   queue := make([]int, 0, d.N)
   d.level[s] = 0
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, e := range d.G[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               queue = append(queue, e.to)
           }
       }
   }
}

func (d *Dinic) dfs(u, t int, f int64) int64 {
   if u == t {
       return f
   }
   for i := d.iter[u]; i < len(d.G[u]); i++ {
       e := &d.G[u][i]
       if e.cap > 0 && d.level[u] < d.level[e.to] {
           ret := d.dfs(e.to, t, min64(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.G[e.to][e.rev].cap += ret
               return ret
           }
       }
       d.iter[u]++
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
   flow := int64(0)
   for {
       d.bfs(s)
       if d.level[t] < 0 {
           break
       }
       for i := range d.iter {
           d.iter[i] = 0
       }
       for {
           f := d.dfs(s, t, INF64)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   return flow
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // compute all-pairs shortest paths (BFS from each node)
   dist := make([][]int, n+1)
   for u := 1; u <= n; u++ {
       d := make([]int, n+1)
       for i := range d {
           d[i] = n + 1
       }
       d[u] = 0
       queue := []int{u}
       for qi := 0; qi < len(queue); qi++ {
           x := queue[qi]
           for _, y := range adj[x] {
               if d[y] > d[x]+1 {
                   d[y] = d[x] + 1
                   queue = append(queue, y)
               }
           }
       }
       dist[u] = d
   }
   var s, b, k int
   fmt.Fscan(in, &s, &b, &k)
   ships := make([]struct{ x int; a, f, p int64 }, s)
   for i := 0; i < s; i++ {
       fmt.Fscan(in, &ships[i].x, &ships[i].a, &ships[i].f, &ships[i].p)
   }
   // bases per planet
   baseDefs := make([][]int64, n+1)
   baseGold := make([][]int64, n+1)
   for i := 0; i < b; i++ {
       var x int
       var d, g int64
       fmt.Fscan(in, &x, &d, &g)
       baseDefs[x] = append(baseDefs[x], d)
       baseGold[x] = append(baseGold[x], g)
   }
   // sort bases per planet by defense and build prefix max of gold
   for v := 1; v <= n; v++ {
       cnt := len(baseDefs[v])
       if cnt == 0 {
           continue
       }
       idx := make([]int, cnt)
       for i := 0; i < cnt; i++ {
           idx[i] = i
       }
       sort.Slice(idx, func(i, j int) bool {
           return baseDefs[v][idx[i]] < baseDefs[v][idx[j]]
       })
       defs := make([]int64, cnt)
       gold := make([]int64, cnt)
       for i, id := range idx {
           defs[i] = baseDefs[v][id]
           gold[i] = baseGold[v][id]
       }
       // prefix max gold
       for i := 1; i < cnt; i++ {
           if gold[i-1] > gold[i] {
               gold[i] = gold[i-1]
           }
       }
       baseDefs[v] = defs
       baseGold[v] = gold
   }
   // build flow network
   N := s + 2
   source := 0
   sink := 1
   dnc := NewDinic(N)
   var sumPos int64
   profits := make([]int64, s)
   for i := 0; i < s; i++ {
       u := ships[i].x
       var best int64 = -INF64
       for v := 1; v <= n; v++ {
           if ships[i].f < int64(dist[u][v]) {
               continue
           }
           defs := baseDefs[v]
           if len(defs) == 0 {
               continue
           }
           // find last index where def <= a
           a := ships[i].a
           idx := sort.Search(len(defs), func(j int) bool { return defs[j] > a })
           if idx > 0 {
               g := baseGold[v][idx-1]
               profit := g - ships[i].p
               if profit > best {
                   best = profit
               }
           }
       }
       profits[i] = best
       node := i + 2
       if best >= 0 {
           sumPos += best
           dnc.AddEdge(source, node, best)
       } else {
           // including unattackable (best==-INF64) gives infinite cap
           cap := -best
           if cap < 0 || cap > INF64 {
               cap = INF64
           }
           dnc.AddEdge(node, sink, cap)
       }
   }
   // dependencies
   for i := 0; i < k; i++ {
       var s1, s2 int
       fmt.Fscan(in, &s1, &s2)
       u := s1 + 1
       v := s2 + 1
       // if u is selected then v must be selected: edge u->v infinite cap
       dnc.AddEdge(u, v, INF64)
   }
   // compute max flow
   flow := dnc.MaxFlow(source, sink)
   res := sumPos - flow
   fmt.Fprintln(out, res)
}

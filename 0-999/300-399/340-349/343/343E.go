package main

import (
   "bufio"
   "fmt"
   "os"
)

// Dinic max-flow implementation
type Edge struct { to, rev, cap int }
// Dinic holds the flow network
type Dinic struct {
   n     int
   adj   [][]Edge
   level []int
   ptr   []int
   q     []int
}

// NewDinic creates a Dinic for n nodes (1-indexed assumed)
func NewDinic(n int) *Dinic {
   adj := make([][]Edge, n+1)
   level := make([]int, n+1)
   ptr := make([]int, n+1)
   q := make([]int, n+1)
   return &Dinic{n, adj, level, ptr, q}
}

// AddEdge adds a directed edge u->v with capacity cap
func (d *Dinic) AddEdge(u, v, cap int) {
   d.adj[u] = append(d.adj[u], Edge{v, len(d.adj[v]), cap})
   d.adj[v] = append(d.adj[v], Edge{u, len(d.adj[u]) - 1, 0})
}

// bfs builds level graph
func (d *Dinic) bfs(s, t int) bool {
   for i := 1; i <= d.n; i++ {
       d.level[i] = -1
   }
   head, tail := 0, 0
   d.q[tail] = s
   d.level[s] = 0
   tail++
   for head < tail {
       u := d.q[head]
       head++
       for _, e := range d.adj[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               d.q[tail] = e.to
               tail++
           }
       }
   }
   return d.level[t] >= 0
}

// dfs sends flow
func (d *Dinic) dfs(u, t, f int) int {
   if u == t || f == 0 {
       return f
   }
   for d.ptr[u] < len(d.adj[u]) {
       e := &d.adj[u][d.ptr[u]]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           pushed := d.dfs(e.to, t, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.adj[e.to][e.rev].cap += pushed
               return pushed
           }
       }
       d.ptr[u]++
   }
   return 0
}

// MaxFlow computes max flow from s to t
func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := 1; i <= d.n; i++ {
           d.ptr[i] = 0
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

// Reachable returns slice of booleans reachable from s in residual graph
func (d *Dinic) Reachable(s int) []bool {
   vis := make([]bool, d.n+1)
   q := make([]int, d.n+1)
   head, tail := 0, 0
   q[tail] = s
   tail++
   vis[s] = true
   for head < tail {
       u := q[head]; head++
       for _, e := range d.adj[u] {
           if e.cap > 0 && !vis[e.to] {
               vis[e.to] = true
               q[tail] = e.to
               tail++
           }
       }
   }
   return vis
}

// Tree edge for permutation building
type TreeEdge struct { to, w, rev, v, next int }

var (
   n, m int
   pipes [][]int
   fa, fv []int
   treeH []int
   treeE []TreeEdge
   treeTot int
   perm []int
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// add tree edge u-v with weight w
func addTreeEdge(u, v, w int) {
   treeE = append(treeE, TreeEdge{v, w, treeTot + 1, 0, treeH[u]})
   treeH[u] = treeTot
   treeTot++
   treeE = append(treeE, TreeEdge{u, w, treeTot - 1, 0, treeH[v]})
   treeH[v] = treeTot
   treeTot++
}

// solve subtree rooted at x, build perm
func solve(x, p int) {
   // find minimum weight unused edge in component
   now := -1
   var dfs func(u, parent int)
   dfs = func(u, parent int) {
       for i := treeH[u]; i != -1; i = treeE[i].next {
           e := &treeE[i]
           if e.to != parent && e.v == 0 {
               if now < 0 || treeE[now].w > e.w {
                   now = i
               }
               dfs(e.to, u)
           }
       }
   }
   dfs(x, p)
   if now < 0 {
       perm = append(perm, x)
       return
   }
   k := now
   treeE[k].v = 1
   treeE[treeE[k].rev].v = 1
   // recurse on two parts
   solve(treeE[k].to, x)
   solve(treeE[treeE[k].rev].to, x)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   pipes = make([][]int, m)
   for i := 0; i < m; i++ {
       pipes[i] = make([]int, 3)
       fmt.Fscan(reader, &pipes[i][0], &pipes[i][1], &pipes[i][2])
   }
   fa = make([]int, n+1)
   fv = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fa[i] = 1
   }
   total := 0
   // Gomory-Hu tree parent array fa, weights fv
   for i := 2; i <= n; i++ {
       d := NewDinic(n)
       // build network for this flow
       for _, e := range pipes {
           u, v, c := e[0], e[1], e[2]
           d.AddEdge(u, v, c)
           d.AddEdge(v, u, c)
       }
       // compute min cut between i and fa[i]
       flow := d.MaxFlow(i, fa[i])
       fv[i] = flow
       total += flow
       // mark reachable from i
       vis := d.Reachable(i)
       for j := i + 1; j <= n; j++ {
           if vis[j] && fa[j] == fa[i] {
               fa[j] = i
           }
       }
   }
   // build tree
   treeH = make([]int, n+1)
   for i := 1; i <= n; i++ {
       treeH[i] = -1
   }
   treeE = make([]TreeEdge, 0, 2*(n-1))
   treeTot = 0
   for i := 2; i <= n; i++ {
       addTreeEdge(fa[i], i, fv[i])
   }
   // get permutation
   perm = make([]int, 0, n)
   solve(1, 0)
   // output
   fmt.Fprintln(writer, total)
   for i, v := range perm {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

// Dinic maxflow implementation
type edge struct { to, rev, cap int }
type Dinic struct {
   N    int
   G    [][]edge
   level []int
   it    []int
}
// NewDinic creates a Dinic for N nodes
func NewDinic(N int) *Dinic {
   return &Dinic{N: N, G: make([][]edge, N), level: make([]int, N), it: make([]int, N)}
}
// Add edge u->v with capacity c
func (d *Dinic) AddEdge(u, v, c int) {
   d.G[u] = append(d.G[u], edge{v, len(d.G[v]), c})
   d.G[v] = append(d.G[v], edge{u, len(d.G[u]) - 1, 0})
}
func min(a, b int) int { if a < b { return a } else { return b } }
func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level { d.level[i] = -1 }
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
   return d.level[t] >= 0
}
func (d *Dinic) dfs(u, t, f int) int {
   if u == t { return f }
   for i := d.it[u]; i < len(d.G[u]); i++ {
       e := &d.G[u][i]
       if e.cap > 0 && d.level[u] < d.level[e.to] {
           ret := d.dfs(e.to, t, min(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.G[e.to][e.rev].cap += ret
               return ret
           }
       }
       d.it[u]++
   }
   return 0
}
// MaxFlow from s to t
func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := range d.it { d.it[i] = 0 }
       for {
           f := d.dfs(s, t, 1e9)
           if f == 0 { break }
           flow += f
       }
   }
   return flow
}

// roundMatrix does matrix rounding for one side
func roundMatrix(d []int, ksum []int, t int) [][]int {
   n := len(d)
   // lower and upper bounds
   lower := make([][]int, n)
   upper := make([][]int, n)
   for i := 0; i < n; i++ {
       fl := d[i] / t
       ce := fl
       if d[i]%t != 0 { ce = fl + 1 }
       lower[i] = make([]int, t)
       upper[i] = make([]int, t)
       for j := 0; j < t; j++ {
           lower[i][j] = fl
           upper[i][j] = ce
       }
   }
   U := n
   // build flow with lower bounds
   // nodes: S=0, u=1..U, j=U+1..U+t, T=U+t+1, SS, TT later
   S := 0
   T := U + t + 1
   N := U + t + 2
   // track demands
   demand := make([]int, N)
   dinic := NewDinic(N + 2)
   SS := N
   TT := N + 1
   // edges S->u cap d[i]
   for i := 0; i < U; i++ {
       dinic.AddEdge(S, 1+i, d[i])
   }
   // edges j->T cap ksum[j]
   for j := 0; j < t; j++ {
       dinic.AddEdge(1+U+j, T, ksum[j])
   }
   // edges u->j with lower and upper
   for i := 0; i < U; i++ {
       for j := 0; j < t; j++ {
           l := lower[i][j]
           ucap := upper[i][j]
           // transform lower bounds: cap = u-l
           dinic.AddEdge(1+i, 1+U+j, ucap-l)
           demand[1+i] -= l
           demand[1+U+j] += l
       }
   }
   // connect T->S infinite cap
   dinic.AddEdge(T, S, 1e9)
   // add edges from SS and to TT for demands
   totalDemand := 0
   for v := 0; v < N; v++ {
       if demand[v] > 0 {
           dinic.AddEdge(SS, v, demand[v])
           totalDemand += demand[v]
       } else if demand[v] < 0 {
           dinic.AddEdge(v, TT, -demand[v])
       }
   }
   // maxflow SS->TT
   if dinic.MaxFlow(SS, TT) != totalDemand {
       // should not happen
       panic("roundMatrix infeasible")
   }
   // now get flows
   A := make([][]int, n)
   for i := range A {
       A[i] = make([]int, t)
   }
   // scan edges from u->j in original list order
   for i := 0; i < U; i++ {
       idx := 0
       for ei := range dinic.G[1+i] {
           e := dinic.G[1+i][ei]
           // to j-node?
           if e.to >= 1+U && e.to < 1+U+t {
               j := e.to - (1 + U)
               l := lower[i][j]
               used := upper[i][j] - l - e.cap
               A[i][j] = l + used
               idx++
           }
       }
   }
   return A
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, k, t int
   fmt.Fscan(reader, &n, &m, &k, &t)
   xs := make([]int, k)
   ys := make([]int, k)
   degL := make([]int, n)
   degR := make([]int, m)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
       xs[i]--
       ys[i]--
       degL[xs[i]]++
       degR[ys[i]]++
   }
   // compute total per company quotas
   ksum := make([]int, t)
   base := k / t
   rem := k % t
   for j := 0; j < t; j++ {
       ksum[j] = base
       if j < rem { ksum[j]++ }
   }
   // round for left and right
   A := roundMatrix(degL, ksum, t)
   B := roundMatrix(degR, ksum, t)
   // prepare edges per left node
   adj := make([][]int, n)
   for i := 0; i < k; i++ {
       adj[xs[i]] = append(adj[xs[i]], i)
   }
   // track unassigned
   col := make([]int, k)
   used := make([]bool, k)
   // assign per company
   for j := 0; j < t; j++ {
       // build flow graph
       N := 2 + n + m
       S := n + m
       T := n + m + 1
       din := NewDinic(N)
       // track each candidate edge and its position in graph
       type emap struct{ u, pos, idx int }
       emaps := make([]emap, 0, k)
       // add source edges and bipartite edges
       for u := 0; u < n; u++ {
           if A[u][j] > 0 {
               din.AddEdge(S, u, A[u][j])
           }
           for _, ei := range adj[u] {
               if used[ei] { continue }
               v := ys[ei]
               // record position
               pos := len(din.G[u])
               din.AddEdge(u, n+v, 1)
               emaps = append(emaps, emap{u, pos, ei})
           }
       }
       for v := 0; v < m; v++ {
           if B[v][j] > 0 {
               din.AddEdge(n+v, T, B[v][j])
           }
       }
       din.MaxFlow(S, T)
       // assign edges where flow was used (cap == 0)
       for _, em := range emaps {
           e := din.G[em.u][em.pos]
           if e.cap == 0 && !used[em.idx] {
               col[em.idx] = j + 1
               used[em.idx] = true
           }
       }
   }
   // compute and output unevenness w_i = max_j cnt - min_j cnt over all vertices
   uneven := 0
   // left side counts
   cntL := make([][]int, n)
   for i := 0; i < n; i++ { cntL[i] = make([]int, t) }
   cntR := make([][]int, m)
   for i := 0; i < m; i++ { cntR[i] = make([]int, t) }
   for i := 0; i < k; i++ {
       cj := col[i] - 1
       cntL[xs[i]][cj]++
       cntR[ys[i]][cj]++
   }
   for i := 0; i < n; i++ {
       mn, mx := cntL[i][0], cntL[i][0]
       for j := 1; j < t; j++ {
           if cntL[i][j] < mn { mn = cntL[i][j] }
           if cntL[i][j] > mx { mx = cntL[i][j] }
       }
       if mx-mn > uneven { uneven = mx - mn }
   }
   for i := 0; i < m; i++ {
       mn, mx := cntR[i][0], cntR[i][0]
       for j := 1; j < t; j++ {
           if cntR[i][j] < mn { mn = cntR[i][j] }
           if cntR[i][j] > mx { mx = cntR[i][j] }
       }
       if mx-mn > uneven { uneven = mx - mn }
   }
   fmt.Fprintln(writer, uneven)
   for i := 0; i < k; i++ {
       fmt.Fprint(writer, col[i], " ")
   }
   fmt.Fprintln(writer)
}

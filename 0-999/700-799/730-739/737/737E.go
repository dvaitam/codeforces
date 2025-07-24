package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Dinic max flow
type edge struct{ to, rev int; cap int }
type Dinic struct{ n, s, t int; g [][]edge; level, it []int }
func newDinic(n, s, t int) *Dinic {
   g := make([][]edge, n)
   return &Dinic{n: n, s: s, t: t, g: g, level: make([]int, n), it: make([]int, n)}
}
func (d *Dinic) addEdge(u, v, c int) {
   d.g[u] = append(d.g[u], edge{v, len(d.g[v]), c})
   d.g[v] = append(d.g[v], edge{u, len(d.g[u]) - 1, 0})
}
func (d *Dinic) bfs() bool {
   for i := range d.level { d.level[i] = -1 }
   q := make([]int, 0, d.n)
   d.level[d.s] = 0; q = append(q, d.s)
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, e := range d.g[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               q = append(q, e.to)
           }
       }
   }
   return d.level[d.t] >= 0
}
func (d *Dinic) dfs(u, f int) int {
   if u == d.t { return f }
   for i := d.it[u]; i < len(d.g[u]); i++ {
       e := &d.g[u][i]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           ret := d.dfs(e.to, min(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.g[e.to][e.rev].cap += ret
               return ret
           }
       }
       d.it[u]++
   }
   return 0
}
func (d *Dinic) maxFlow() int {
   flow := 0
   for d.bfs() {
       for i := range d.it { d.it[i] = 0 }
       for {
           f := d.dfs(d.s, 1e9)
           if f == 0 { break }
           flow += f
       }
   }
   return flow
}
func min(a, b int) int { if a < b { return a }; return b }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, B int
   fmt.Fscan(in, &n, &m, &B)
   p := make([]int, m)
   for i := 0; i < m; i++ { fmt.Fscan(in, &p[i]) }
   // tasks: rem[i][j]
   rem0 := make([][]int, n)
   sumChild := make([]int, n)
   sumMach := make([]int, m)
   total := 0
   for i := 0; i < n; i++ {
       var k int; fmt.Fscan(in, &k)
       rem0[i] = make([]int, m)
       for z := 0; z < k; z++ {
           var j, t int; fmt.Fscan(in, &j, &t); j--
           rem0[i][j] = t
           sumChild[i] += t
           sumMach[j] += t
           total += t
       }
   }
   bestT := 1<<60; bestMask := 0
   // enumerate masks
   for mask := 0; mask < (1 << m); mask++ {
       cost := 0
       for j := 0; j < m; j++ {
           if (mask>>j)&1 == 1 { cost += p[j] }
       }
       if cost > B { continue }
       // capacity x[j]
       x := make([]int, m)
       for j := 0; j < m; j++ { x[j] = 1 + ((mask>>j)&1) }
       // lower bound for T
       lo := 0
       for i := 0; i < n; i++ { lo = max(lo, sumChild[i]) }
       for j := 0; j < m; j++ {
           need := (sumMach[j] + x[j] - 1) / x[j]
           lo = max(lo, need)
       }
       hi := total
       if lo >= bestT { continue }
       // binary search T
       for lo <= hi {
           mid := (lo + hi) / 2
           if mid >= bestT { hi = mid - 1; continue }
           if feasible(n, m, rem0, x, mid) {
               bestT = mid; bestMask = mask
               hi = mid - 1
           } else {
               lo = mid + 1
           }
       }
   }
   // output bestT and mask
   fmt.Println(bestT)
   mask := bestMask
   xs := make([]int, m)
   for j := 0; j < m; j++ { xs[j] = (mask>>j)&1 }
   for j := 0; j < m; j++ { fmt.Print(xs[j]) }
   fmt.Println()
   // prepare rem copy and simulate
   rem := make([][]int, n)
   for i := range rem0 {
       rem[i] = make([]int, m)
       copy(rem[i], rem0[i])
   }
   // machine copies availability lists
   childBusy := make([]bool, n)
   type seg struct{ i, j, s int }
   var segs []seg
   // simulate time 0..bestT-1
   for t := 0; t < bestT; t++ {
       for i := range childBusy { childBusy[i] = false }
       for j := 0; j < m; j++ {
           copies := 1 + ((bestMask>>j)&1)
           for c := 0; c < copies; c++ {
               // find a child i to run
               for i := 0; i < n; i++ {
                   if rem[i][j] > 0 && !childBusy[i] {
                       rem[i][j]--
                       childBusy[i] = true
                       segs = append(segs, seg{i + 1, j + 1, t})
                       break
                   }
               }
           }
       }
   }
   // output segments
   g := len(segs)
   fmt.Println(g)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, s := range segs {
       fmt.Fprintf(w, "%d %d %d 1\n", s.i, s.j, s.s)
   }
}

func feasible(n, m int, rem0 [][]int, x []int, T int) bool {
   // build flow
   N := 2 + m + n
   s, t := 0, N-1
   d := newDinic(N, s, t)
   total := 0
   for j := 0; j < m; j++ {
       d.addEdge(s, 1+j, T*x[j])
       for i := 0; i < n; i++ {
           if rem0[i][j] > 0 {
               d.addEdge(1+j, 1+m+i, rem0[i][j])
               total += rem0[i][j]
           }
       }
   }
   for i := 0; i < n; i++ {
       d.addEdge(1+m+i, t, T)
   }
   return d.maxFlow() == total
}
func max(a, b int) int { if a > b { return a }; return b }

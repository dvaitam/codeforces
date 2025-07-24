package main

import (
   "bufio"
   "fmt"
   "os"
)

// Max flow implementation (Dinic)
type edge struct{ to, rev, cap int }
type dinic struct{ n, s, t int; g [][]edge; level, ptr []int }
func newDinic(n, s, t int) *dinic {
   g := make([][]edge, n)
   return &dinic{n, s, t, g, make([]int, n), make([]int, n)}
}
func (d *dinic) addEdge(u, v, c int) {
   d.g[u] = append(d.g[u], edge{v, len(d.g[v]), c})
   d.g[v] = append(d.g[v], edge{u, len(d.g[u]) - 1, 0})
}
func (d *dinic) bfs() bool {
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
func (d *dinic) dfs(u, f int) int {
   if u == d.t || f == 0 { return f }
   for ; d.ptr[u] < len(d.g[u]); d.ptr[u]++ {
       e := &d.g[u][d.ptr[u]]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           pushed := d.dfs(e.to, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.g[e.to][e.rev].cap += pushed
               return pushed
           }
       }
   }
   return 0
}
func (d *dinic) MaxFlow() int {
   flow := 0
   for d.bfs() {
       for i := range d.ptr { d.ptr[i] = 0 }
       for pushed := d.dfs(d.s, 1<<60); pushed > 0; pushed = d.dfs(d.s, 1<<60) {
           flow += pushed
       }
   }
   return flow
}
func min(a, b int) int { if a < b { return a }; return b }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // compute groups by digit length
   pw := []int{1}
   for i := 1; i <= 10; i++ { pw = append(pw, pw[i-1]*10) }
   D := 0
   for D = 1; pw[D] <= n; D++ {
   }
   // groups 1..D
   c := make([]int, D+1)
   groupNodes := make([][]int, D+1)
   for k := 1; k <= D; k++ {
       l := pw[k-1]
       r := min(n, pw[k]-1)
       if r >= l {
           for x := l; x <= r; x++ {
               groupNodes[k] = append(groupNodes[k], x)
           }
           c[k] = r - l + 1
       }
   }
   // read blueprint edges
   type pair struct{ i, j int }
   edges := make([]pair, n-1)
   dcnt := make([][]int, D+1)
   for i := range dcnt { dcnt[i] = make([]int, D+1) }
   for idx := 0; idx < n-1; idx++ {
       var a, b string
       fmt.Fscan(in, &a, &b)
       i := len(a)
       j := len(b)
       if i < 1 || i > D || j < 1 || j > D {
           fmt.Println(-1); return
       }
       edges[idx] = pair{i, j}
       dcnt[i][j]++
       if i != j { dcnt[j][i]++ }
   }
   // check H' connectivity on groups with c[k]>0
   usedGrp := make([]bool, D+1)
   for k := 1; k <= D; k++ {
       if c[k] > 0 {
           usedGrp[k] = true
       }
   }
   seen := make([]bool, D+1)
   var dfsGrp func(int)
   dfsGrp = func(u int) {
       seen[u] = true
       for v := 1; v <= D; v++ {
           if !seen[v] && u != v && dcnt[u][v] > 0 {
               dfsGrp(v)
           }
       }
   }
   // find first group as root
   root := 0
   for k := 1; k <= D; k++ {
       if usedGrp[k] { root = k; break }
   }
   if root == 0 {
       fmt.Println(-1); return
   }
   dfsGrp(root)
   for k := 1; k <= D; k++ {
       if usedGrp[k] && !seen[k] {
           fmt.Println(-1); return
       }
   }
   // build spanning tree on groups H'
   parent := make([]int, D+1)
   for i := range parent { parent[i] = -1 }
   parent[root] = 0
   queue := []int{root}
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for v := 1; v <= D; v++ {
           if parent[v] < 0 && u != v && dcnt[u][v] > 0 {
               parent[v] = u
               queue = append(queue, v)
           }
       }
   }
   // backbone edges list of types
   back := make(map[pair]int)
   for v := 1; v <= D; v++ {
       u := parent[v]
       if u > 0 {
           // use one edge u-v
           back[pair{u, v}] = 1
           back[pair{v, u}] = 1
           // decrement dcnt for spanning tree selection
           dcnt[u][v]--
           dcnt[v][u]--
       }
   }
   // leftover x_ij counts for i<=j
   type tkey struct{ i, j int }
   xmap := make(map[tkey]int)
   totalLeft := 0
   for i := 1; i <= D; i++ {
       for j := i; j <= D; j++ {
           cnt := dcnt[i][j]
           if i != j {
               cnt = dcnt[i][j] // after decrement both sides
           }
           if cnt > 0 {
               xmap[tkey{i, j}] = cnt
               totalLeft += cnt
           }
       }
   }
   // leaves per group
   leaves := make([]int, D+1)
   for k := 1; k <= D; k++ {
       if c[k] > 0 {
           leaves[k] = c[k] - 1
       }
   }
   // build flow network
   // nodes: source 0, types 1..T, groups T+1..T+D, sink T+D+1
   typeKeys := make([]tkey, 0, len(xmap))
   for k := range xmap { typeKeys = append(typeKeys, k) }
   Tn := len(typeKeys)
   N := 1 + Tn + D + 1
   S, Tt := 0, N-1
   din := newDinic(N, S, Tt)
   // map type index
   tidx := make(map[tkey]int)
   for idx, k := range typeKeys {
       tidx[k] = idx + 1
       din.addEdge(S, idx+1, xmap[k])
       // to groups
       i, j := k.i, k.j
       if i == j {
           din.addEdge(idx+1, Tn+i, xmap[k])
       } else {
           din.addEdge(idx+1, Tn+i, xmap[k])
           din.addEdge(idx+1, Tn+j, xmap[k])
       }
   }
   sumLeaves := 0
   for k := 1; k <= D; k++ {
       din.addEdge(Tn+k, Tt, leaves[k])
       sumLeaves += leaves[k]
   }
   if sumLeaves != totalLeft || din.MaxFlow() != totalLeft {
       fmt.Println(-1); return
   }
   // collect flow from types to groups
   // f_e_k = initial cap - residual cap on edge idx
   // Build mapping of type index to flow to i
   flowUse := make(map[tkey]int)
   for idx, k := range typeKeys {
       u := idx + 1
       for _, e := range din.g[u] {
           // edges to group nodes
           if e.to > Tn && e.to < Tt {
               grp := e.to - Tn
               used := xmap[k] - e.cap
               if grp == k.i {
                   flowUse[k] += used
               } else if grp == k.j {
                   // consumed on j side
               }
           }
       }
   }
   // assign actual edges
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // pick centers
   center := make([]int, D+1)
   for k := 1; k <= D; k++ {
       if len(groupNodes[k]) > 0 {
           center[k] = groupNodes[k][0]
           groupNodes[k] = groupNodes[k][1:]
       }
   }
   // output backbone edges
   for v := 1; v <= D; v++ {
       u := parent[v]
       if u > 0 {
           fmt.Fprintf(out, "%d %d\n", center[u], center[v])
       }
   }
   // output leftover edges by types
   for _, k := range typeKeys {
       cnt := xmap[k]
       usedI := flowUse[k]
       usedJ := cnt - usedI
       i, j := k.i, k.j
       // loops, j==i, usedI==cnt
       // first useI edges leaf i -> center j
       for t := 0; t < usedI; t++ {
           leaf := groupNodes[i][len(groupNodes[i])-1]
           groupNodes[i] = groupNodes[i][:len(groupNodes[i])-1]
           fmt.Fprintf(out, "%d %d\n", leaf, center[j])
       }
       // then usedJ edges leaf j -> center i (if i!=j)
       for t := 0; t < usedJ; t++ {
           leaf := groupNodes[j][len(groupNodes[j])-1]
           groupNodes[j] = groupNodes[j][:len(groupNodes[j])-1]
           fmt.Fprintf(out, "%d %d\n", leaf, center[i])
       }
   }
}

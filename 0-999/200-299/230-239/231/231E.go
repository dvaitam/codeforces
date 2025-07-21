package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const MOD = 1000000007

type Edge struct{ to, id int }
// DFS frame for cycle detection
type frame struct{ u, peid, dep, idx int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // read n, m
   var n, m int
   fmt.Fscan(in, &n, &m)
   g := make([][]Edge, n)
   uE := make([]int, m)
   vE := make([]int, m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--, v--
       uE[i], vE[i] = u, v
       g[u] = append(g[u], Edge{v, i})
       g[v] = append(g[v], Edge{u, i})
   }
   // cycle detection
   visited := make([]int8, n) // 0,1,2
   parent := make([]int, n)
   parentE := make([]int, n)
   depth := make([]int, n)
   edgeInCycle := make([]bool, m)
   var cycles [][]int
   // iterative DFS
   var stk []frame
   stk = append(stk, frame{u: 0, peid: -1, dep: 0, idx: 0})
   parent[0] = -1; parentE[0] = -1; depth[0] = 0
   for len(stk) > 0 {
       f := &stk[len(stk)-1]
       u := f.u
       if visited[u] == 0 {
           visited[u] = 1
           parentE[u] = f.peid
           depth[u] = f.dep
       }
       if f.idx < len(g[u]) {
           e := g[u][f.idx]; f.idx++
           v := e.to; eid := e.id
           if eid == parentE[u] {
               continue
           }
           if visited[v] == 0 {
               parent[v] = u
               stk = append(stk, frame{u: v, peid: eid, dep: f.dep + 1, idx: 0})
               continue
           }
           if visited[v] == 1 && depth[v] < f.dep {
               // found cycle u->v
               cid := len(cycles)
               var cyc []int
               // mark this edge
               edgeInCycle[eid] = true
               w := u
               for w != v {
                   cyc = append(cyc, w)
                   edgeInCycle[parentE[w]] = true
                   w = parent[w]
               }
               cyc = append(cyc, v)
               cycles = append(cycles, cyc)
           }
       } else {
           visited[u] = 2
           stk = stk[:len(stk)-1]
       }
   }
   // build block tree
   cCnt := len(cycles)
   tot := n + cCnt
   bAdj := make([][]int, tot)
   // add tree edges (non-cycle)
   for i := 0; i < m; i++ {
       if edgeInCycle[i] {
           continue
       }
       u, v := uE[i], vE[i]
       bAdj[u] = append(bAdj[u], v)
       bAdj[v] = append(bAdj[v], u)
   }
   // add cycle block nodes
   isCycleNode := make([]bool, tot)
   for cid, cyc := range cycles {
       bid := n + cid
       isCycleNode[bid] = true
       for _, u := range cyc {
           bAdj[bid] = append(bAdj[bid], u)
           bAdj[u] = append(bAdj[u], bid)
       }
   }
   // LCA prep on block tree
   LOG := 1
   for (1 << LOG) <= tot { LOG++ }
   up := make([][]int, LOG)
   for i := range up { up[i] = make([]int, tot) }
   depthB := make([]int, tot)
   distC := make([]int, tot)
   // BFS from 0
   queue := make([]int, 0, tot)
   parentB := make([]int, tot)
   queue = append(queue, 0)
   parentB[0] = 0
   depthB[0] = 0
   distC[0] = 0
   up[0][0] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range bAdj[u] {
           if v == parentB[u] {
               continue
           }
           parentB[v] = u
           depthB[v] = depthB[u] + 1
           distC[v] = distC[u]
           if isCycleNode[v] {
               distC[v]++
           }
           up[0][v] = u
           queue = append(queue, v)
       }
   }
   for k1 := 1; k1 < LOG; k1++ {
       for v := 0; v < tot; v++ {
           up[k1][v] = up[k1-1][ up[k1-1][v] ]
       }
   }
   // precompute powers of 2
   maxC := cCnt
   pow2 := make([]int, maxC+2)
   pow2[0] = 1
   for i := 1; i <= maxC; i++ {
       pow2[i] = pow2[i-1] * 2 % MOD
   }
   // queries
   var k int
   fmt.Fscan(in, &k)
   for i := 0; i < k; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--, y--
       l := lca(x, y, depthB, up, LOG)
       cnt := distC[x] + distC[y] - 2*distC[l]
       if isCycleNode[l] {
           cnt++
       }
       if cnt < 0 {
           cnt = 0
       }
       if cnt <= maxC {
           fmt.Fprintln(out, pow2[cnt])
       } else {
           // fallback
           res := 1
           for j := 0; j < cnt; j++ {
               res = res * 2 % MOD
           }
           fmt.Fprintln(out, res)
       }
   }
}

func lca(a, b int, depth []int, up [][]int, LOG int) int {
   if depth[a] < depth[b] {
       a, b = b, a
   }
   diff := depth[a] - depth[b]
   for k := 0; diff > 0; k++ {
       if diff&1 != 0 {
           a = up[k][a]
       }
       diff >>= 1
   }
   if a == b {
       return a
   }
   for k := LOG - 1; k >= 0; k-- {
       ua := up[k][a]; ub := up[k][b]
       if ua != ub {
           a = ua; b = ub
       }
   }
   return up[0][a]
}

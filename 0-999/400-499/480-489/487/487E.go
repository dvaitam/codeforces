package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1<<60

type edge struct{u, v int}

var (
   n, m, q int
   wgt []int64
   adj [][]int
   // for biconnected
   dfn, low []int
   timer int
   stk []edge
   mark []int
   curMark, bcCnt int
   // block tree
   btAdj [][]int
   totN int
   // HLD
   parent, depth, heavy, head, pos, sz []int
   curPos int
   flatW []int64
   // segtree
   segN int
   segT []int64
)

func minInt(a, b int) int {
   if a < b { return a }
   return b
}
func min64(a, b int64) int64 {
   if a < b { return a }
   return b
}

func tarjan(u, pe int) {
   timer++
   dfn[u], low[u] = timer, timer
   for _, v := range adj[u] {
       if v == pe { pe = -1; continue }
       if dfn[v] == 0 {
           stk = append(stk, edge{u, v})
           tarjan(v, u)
           low[u] = minInt(low[u], low[v])
           if low[v] >= dfn[u] {
               bcCnt++
               bid := n + bcCnt
               curMark++
               var e edge
               for {
                   e, stk = stk[len(stk)-1], stk[:len(stk)-1]
                   for _, x := range []int{e.u, e.v} {
                       if mark[x] != curMark {
                           mark[x] = curMark
                           // connect block node and vertex
                           btAdj[bid] = append(btAdj[bid], x)
                           btAdj[x] = append(btAdj[x], bid)
                       }
                   }
                   if e.u == u && e.v == v {
                       break
                   }
               }
           }
       } else if dfn[v] < dfn[u] {
           // back edge
           stk = append(stk, edge{u, v})
           low[u] = minInt(low[u], dfn[v])
       }
   }
}

func dfs1(u, p int) {
   sz[u] = 1
   parent[u] = p
   depth[u] = depth[p] + 1
   maxSz := 0
   for _, v := range btAdj[u] {
       if v == p { continue }
       dfs1(v, u)
       if sz[v] > maxSz {
           maxSz = sz[v]
           heavy[u] = v
       }
       sz[u] += sz[v]
   }
}

func dfs2(u, h int) {
   head[u] = h
   pos[u] = curPos
   // set flat weight
   if u <= n {
       flatW[curPos] = wgt[u]
   } else {
       flatW[curPos] = INF
   }
   curPos++
   if heavy[u] != 0 {
       dfs2(heavy[u], h)
   }
   for _, v := range btAdj[u] {
       if v == parent[u] || v == heavy[u] { continue }
       dfs2(v, v)
   }
}

// segtree on [1..totN]
func segInit(n0 int) {
   segN = 1
   for segN < n0 { segN <<= 1 }
   segT = make([]int64, 2*segN)
   for i := 1; i < 2*segN; i++ {
       segT[i] = INF
   }
   for i := 1; i <= n0; i++ {
       segT[segN+i-1] = flatW[i]
   }
   for i := segN-1; i >= 1; i-- {
       segT[i] = min64(segT[2*i], segT[2*i+1])
   }
}
func segUpdate(p int, v int64) {
   i := segN + p - 1
   segT[i] = v
   for i >>= 1; i > 0; i >>= 1 {
       segT[i] = min64(segT[2*i], segT[2*i+1])
   }
}
func segQuery(l, r int) int64 {
   res := INF
   l += segN - 1; r += segN - 1
   for l <= r {
       if l&1 == 1 {
           res = min64(res, segT[l]); l++
       }
       if r&1 == 0 {
           res = min64(res, segT[r]); r--
       }
       l >>= 1; r >>= 1
   }
   return res
}

func pathQuery(u, v int) int64 {
   res := int64(INF)
   for head[u] != head[v] {
       if depth[head[u]] > depth[head[v]] {
           res = min64(res, segQuery(pos[head[u]], pos[u]))
           u = parent[head[u]]
       } else {
           res = min64(res, segQuery(pos[head[v]], pos[v]))
           v = parent[head[v]]
       }
   }
   l, r := pos[u], pos[v]
   if l > r { l, r = r, l }
   res = min64(res, segQuery(l, r))
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m, &q)
   wgt = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &wgt[i])
   }
   adj = make([][]int, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   // init tarjan & block tree
   dfn = make([]int, n+1)
   low = make([]int, n+1)
   mark = make([]int, n+1)
   btAdj = make([][]int, n+m+5)
   // run tarjan
   for i := 1; i <= n; i++ {
       if dfn[i] == 0 {
           tarjan(i, -1)
       }
   }
   totN = n + bcCnt
   // HLD
   parent = make([]int, totN+1)
   depth = make([]int, totN+1)
   heavy = make([]int, totN+1)
   sz = make([]int, totN+1)
   dfs1(1, 0)
   head = make([]int, totN+1)
   pos = make([]int, totN+1)
   flatW = make([]int64, totN+1)
   curPos = 1
   dfs2(1, 1)
   // build segtree
   segInit(totN)
   // process queries
   for i := 0; i < q; i++ {
       var typ byte
       fmt.Fscan(in, &typ)
       if typ == 'C' {
           var a int; var v int64
           fmt.Fscan(in, &a, &v)
           wgt[a] = v
           segUpdate(pos[a], v)
       } else if typ == 'A' {
           var a, b int
           fmt.Fscan(in, &a, &b)
           ans := pathQuery(a, b)
           fmt.Fprintln(out, ans)
       }
   }
}

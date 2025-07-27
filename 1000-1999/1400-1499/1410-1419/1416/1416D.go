package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU with merge tree
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   a := make([]int, m+1)
   b := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   // queries
   typ := make([]int, q+1)
   arg := make([]int, q+1)
   dt := make([]int, m+1)
   for i := 1; i <= m; i++ {
       dt[i] = q + 1
   }
   for i := 1; i <= q; i++ {
       fmt.Fscan(reader, &typ[i], &arg[i])
       if typ[i] == 2 {
           ei := arg[i]
           dt[ei] = i
       }
   }
   // edges sorted by deletion time descending
   type edgeInfo struct{ t, idx int }
   edges := make([]edgeInfo, m)
   for i := 1; i <= m; i++ {
       edges[i-1] = edgeInfo{dt[i], i}
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].t > edges[j].t
   })
   // DSU
   maxNodes := 2*n + 5
   dsu := make([]int, maxNodes)
   parent := make([]int, maxNodes)
   mergeTime := make([]int, maxNodes)
   left := make([]int, maxNodes)
   right := make([]int, maxNodes)
   for i := 1; i < maxNodes; i++ {
       dsu[i] = i
       mergeTime[i] = q + 1
   }
   tot := n
   var find func(int) int
   find = func(x int) int {
       if dsu[x] != x {
           dsu[x] = find(dsu[x])
       }
       return dsu[x]
   }
   for _, e := range edges {
       if e.t <= 0 {
           break
       }
       u := find(a[e.idx])
       v := find(b[e.idx])
       if u != v {
           tot++
           mergeTime[tot] = e.t
           left[tot] = u
           right[tot] = v
           parent[u] = tot
           parent[v] = tot
           dsu[u] = tot
           dsu[v] = tot
           dsu[tot] = tot
       }
   }
   // binary lifting
   LOG := 0
   for (1<<LOG) <= tot {
       LOG++
   }
   up := make([][]int, LOG)
   for i := 0; i < LOG; i++ {
       up[i] = make([]int, tot+1)
   }
   for v := 1; v <= tot; v++ {
       up[0][v] = parent[v]
   }
   for i := 1; i < LOG; i++ {
       for v := 1; v <= tot; v++ {
           up[i][v] = up[i-1][ up[i-1][v] ]
       }
   }
   // build tree and dfs order
   g := make([][]int, tot+1)
   for v := 1; v <= tot; v++ {
       if parent[v] != 0 {
           g[parent[v]] = append(g[parent[v]], v)
       }
   }
   tin := make([]int, tot+1)
   tout := make([]int, tot+1)
   timer := 0
   var dfs func(int)
   dfs = func(u int) {
       timer++
       tin[u] = timer
       for _, v := range g[u] {
           dfs(v)
       }
       tout[u] = timer
   }
   // dfs from roots
   for v := 1; v <= tot; v++ {
       if parent[v] == 0 {
           dfs(v)
       }
   }
   // build segment tree
   size := 1
   for size < tot+2 {
       size <<= 1
   }
   type Node struct{ val, idx int }
   seg := make([]Node, 2*size)
   // initialize leaves
   for i := 1; i <= n; i++ {
       pos := tin[i]
       seg[size+pos] = Node{p[i], pos}
   }
   // build
   for i := size - 1; i >= 1; i-- {
       leftN := seg[2*i]
       rightN := seg[2*i+1]
       if leftN.val >= rightN.val {
           seg[i] = leftN
       } else {
           seg[i] = rightN
       }
   }
   // helper to get comp root at time t
   getRoot := func(v, ttime int) int {
       cur := v
       for i := LOG - 1; i >= 0; i-- {
           u := up[i][cur]
           if u != 0 && mergeTime[u] > ttime {
               cur = u
           }
       }
       return cur
   }
   // process queries
   for i := 1; i <= q; i++ {
       if typ[i] == 1 {
           v := arg[i]
           comp := getRoot(v, i)
           l, r := tin[comp], tout[comp]
           // range max query
           l += size; r += size
           res := Node{0, l}
           for l <= r {
               if l&1 == 1 {
                   if seg[l].val > res.val {
                       res = seg[l]
                   }
                   l++
               }
               if r&1 == 0 {
                   if seg[r].val > res.val {
                       res = seg[r]
                   }
                   r--
               }
               l >>= 1; r >>= 1
           }
           // output
           fmt.Fprintln(writer, res.val)
           if res.val > 0 {
               // update to 0 at res.idx
               pos := res.idx + size
               seg[pos].val = 0
               pos >>= 1
               for pos >= 1 {
                   leftN := seg[2*pos]
                   rightN := seg[2*pos+1]
                   if leftN.val >= rightN.val {
                       seg[pos] = leftN
                   } else {
                       seg[pos] = rightN
                   }
                   pos >>= 1
               }
           }
       }
   }
}

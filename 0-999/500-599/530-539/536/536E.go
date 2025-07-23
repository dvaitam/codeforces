package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   to int
   w  int
}

var n, q int
var f []int64
var adj [][]Edge
var parent, depth, heavy, head, pos []int
var weightAtNode []int
var curPos int

// Data holds segment info
type Data struct {
   length, pref, suf int
   sum               int64
}

func combine(a, b Data) Data {
   if a.length == 0 {
       return b
   }
   if b.length == 0 {
       return a
   }
   res := Data{length: a.length + b.length}
   res.sum = a.sum + b.sum
   // merge block if needed
   if a.suf > 0 && b.pref > 0 {
       res.sum -= f[a.suf]
       res.sum -= f[b.pref]
       res.sum += f[a.suf+b.pref]
   }
   // prefix
   res.pref = a.pref
   if a.pref == a.length {
       res.pref = a.length + b.pref
   }
   // suffix
   res.suf = b.suf
   if b.suf == b.length {
       res.suf = b.length + a.suf
   }
   return res
}

func reverseData(a Data) Data {
   return Data{length: a.length, pref: a.suf, suf: a.pref, sum: a.sum}
}

// segment tree
var st []Data

func buildST(n int) {
   st = make([]Data, 4*(n+1))
   // initialize leaves to length 1 (val=0)
   var init func(idx, l, r int)
   init = func(idx, l, r int) {
       if l == r {
           st[idx] = Data{length: 1}
           return
       }
       mid := (l + r) >> 1
       init(idx*2, l, mid)
       init(idx*2+1, mid+1, r)
       st[idx] = combine(st[idx*2], st[idx*2+1])
   }
   init(1, 1, n)
}

func updateST(idx, l, r, posi, val int) {
   if l == r {
       if val == 1 {
           st[idx] = Data{length: 1, pref: 1, suf: 1, sum: f[1]}
       } else {
           st[idx] = Data{length: 1}
       }
       return
   }
   mid := (l + r) >> 1
   if posi <= mid {
       updateST(idx*2, l, mid, posi, val)
   } else {
       updateST(idx*2+1, mid+1, r, posi, val)
   }
   st[idx] = combine(st[idx*2], st[idx*2+1])
}

func queryST(idx, l, r, ql, qr int) Data {
   if ql > r || qr < l {
       return Data{} // zero
   }
   if ql <= l && r <= qr {
       return st[idx]
   }
   mid := (l + r) >> 1
   left := queryST(idx*2, l, mid, ql, qr)
   right := queryST(idx*2+1, mid+1, r, ql, qr)
   return combine(left, right)
}

// dfs to compute sizes and heavy child
func dfs(u, p int) int {
   size := 1
   maxSub := 0
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       parent[v] = u
       depth[v] = depth[u] + 1
       weightAtNode[v] = e.w
       sub := dfs(v, u)
       if sub > maxSub {
           maxSub = sub
           heavy[u] = v
       }
       size += sub
   }
   return size
}

func decompose(u, h int) {
   head[u] = h
   curPos++
   pos[u] = curPos
   if heavy[u] != 0 {
       decompose(heavy[u], h)
   }
   for _, e := range adj[u] {
       v := e.to
       if v == parent[u] || v == heavy[u] {
           continue
       }
       decompose(v, v)
   }
}

func queryPath(u, v int) int64 {
   var parts []Data
   lcaU, lcaV := u, v
   // find lca
   for head[lcaU] != head[lcaV] {
       if depth[head[lcaU]] > depth[head[lcaV]] {
           lcaU = parent[head[lcaU]]
       } else {
           lcaV = parent[head[lcaV]]
       }
   }
   var lca int
   if depth[lcaU] < depth[lcaV] {
       lca = lcaU
   } else {
       lca = lcaV
   }
   // u to lca (u is v original)
   u = u
   for head[u] != head[lca] {
       d := queryST(1, 1, n, pos[head[u]], pos[u])
       parts = append(parts, reverseData(d))
       u = parent[head[u]]
   }
   if u != lca {
       d := queryST(1, 1, n, pos[lca]+1, pos[u])
       parts = append(parts, reverseData(d))
   }
   // lca to v
   var parts2 []Data
   v = v
   for head[v] != head[lca] {
       d := queryST(1, 1, n, pos[head[v]], pos[v])
       parts2 = append(parts2, d)
       v = parent[head[v]]
   }
   if v != lca {
       d := queryST(1, 1, n, pos[lca]+1, pos[v])
       parts2 = append(parts2, d)
   }
   // reverse parts2
   for i, j := 0, len(parts2)-1; i < j; i, j = i+1, j-1 {
       parts2[i], parts2[j] = parts2[j], parts2[i]
   }
   // combine all
   res := Data{}
   for _, d := range parts {
       res = combine(res, d)
   }
   for _, d := range parts2 {
       res = combine(res, d)
   }
   return res.sum
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &q)
   f = make([]int64, n+1)
   for i := 1; i < n; i++ {
       fmt.Fscan(in, &f[i])
   }
   adj = make([][]Edge, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       var w int
       fmt.Fscan(in, &u, &v, &w)
       adj[u] = append(adj[u], Edge{to: v, w: w})
       adj[v] = append(adj[v], Edge{to: u, w: w})
   }
   parent = make([]int, n+1)
   depth = make([]int, n+1)
   heavy = make([]int, n+1)
   head = make([]int, n+1)
   pos = make([]int, n+1)
   weightAtNode = make([]int, n+1)
   depth[1] = 0
   parent[1] = 0
   dfs(1, 0)
   curPos = 0
   decompose(1, 1)
   buildST(n)
   // prepare nodes by weight
   type NodeW struct{ w, u int }
   var nodes []NodeW
   nodes = make([]NodeW, 0, n-1)
   for u := 2; u <= n; u++ {
       nodes = append(nodes, NodeW{w: weightAtNode[u], u: u})
   }
   sort.Slice(nodes, func(i, j int) bool { return nodes[i].w > nodes[j].w })
   // read queries
   type Query struct{ l, u, v, id int }
   qs := make([]Query, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(in, &qs[i].v, &qs[i].u, &qs[i].l)
       qs[i].id = i
   }
   sort.Slice(qs, func(i, j int) bool { return qs[i].l > qs[j].l })
   ans := make([]int64, q)
   idx := 0
   for _, qu := range qs {
       for idx < len(nodes) && nodes[idx].w >= qu.l {
           updateST(1, 1, n, pos[nodes[idx].u], 1)
           idx++
       }
       ans[qu.id] = queryPath(qu.v, qu.u)
   }
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

// Segment tree with range invert and range sum
type SegTree struct {
   n    int
   sum  []int
   lazy []bool
}

func NewSegTree(n int) *SegTree {
   st := &SegTree{n: n}
   st.sum = make([]int, 4*n)
   st.lazy = make([]bool, 4*n)
   return st
}

func (st *SegTree) apply(node, l, r int) {
   st.sum[node] = (r - l + 1) - st.sum[node]
   st.lazy[node] = !st.lazy[node]
}

func (st *SegTree) push(node, l, r int) {
   if st.lazy[node] {
       m := (l + r) >> 1
       st.apply(node*2, l, m)
       st.apply(node*2+1, m+1, r)
       st.lazy[node] = false
   }
}

func (st *SegTree) update(node, l, r, ql, qr int) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       st.apply(node, l, r)
       return
   }
   st.push(node, l, r)
   m := (l + r) >> 1
   st.update(node*2, l, m, ql, qr)
   st.update(node*2+1, m+1, r, ql, qr)
   st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
}

func (st *SegTree) query(node, l, r, ql, qr int) int {
   if ql > r || qr < l {
       return 0
   }
   if ql <= l && r <= qr {
       return st.sum[node]
   }
   st.push(node, l, r)
   m := (l + r) >> 1
   return st.query(node*2, l, m, ql, qr) + st.query(node*2+1, m+1, r, ql, qr)
}

// HLD on forest
var (
   adj       [][]int
   treeAdj   [][]int
   parentArr []int
   depthArr  []int
   heavy     []int
   sizeArr   []int
   head      []int
   pos       []int
   curPos    int
)

func dfsSz(u int) {
   sizeArr[u] = 1
   maxSz := 0
   for _, v := range treeAdj[u] {
       if v == parentArr[u] {
           continue
       }
       parentArr[v] = u
       depthArr[v] = depthArr[u] + 1
       dfsSz(v)
       if sizeArr[v] > maxSz {
           maxSz = sizeArr[v]
           heavy[u] = v
       }
       sizeArr[u] += sizeArr[v]
   }
}

func decompose(u, h int) {
   head[u] = h
   pos[u] = curPos
   curPos++
   if heavy[u] != -1 {
       decompose(heavy[u], h)
   }
   for _, v := range treeAdj[u] {
       if v != parentArr[u] && v != heavy[u] {
           decompose(v, v)
       }
   }
}

// process path u-v: apply f on segments
func pathOp(u, v int, st *SegTree, opType int) int {
   // opType: 0=query sum, 1=update invert
   res := 0
   for head[u] != head[v] {
       if depthArr[head[u]] > depthArr[head[v]] {
           if opType == 0 {
               res += st.query(1, 0, st.n-1, pos[head[u]], pos[u])
           } else {
               st.update(1, 0, st.n-1, pos[head[u]], pos[u])
           }
           u = parentArr[head[u]]
       } else {
           if opType == 0 {
               res += st.query(1, 0, st.n-1, pos[head[v]], pos[v])
           } else {
               st.update(1, 0, st.n-1, pos[head[v]], pos[v])
           }
           v = parentArr[head[v]]
       }
   }
   if u == v {
       return res
   }
   if depthArr[u] > depthArr[v] {
       u, v = v, u
   }
   // u is lca, skip pos[u]
   if opType == 0 {
       res += st.query(1, 0, st.n-1, pos[u]+1, pos[v])
   } else {
       st.update(1, 0, st.n-1, pos[u]+1, pos[v])
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   adj = make([][]int, n+1)
   edges := make([][2]int, n)
   for i := 0; i < n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       edges[i][0], edges[i][1] = u, v
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // find cycle by peeling leaves
   degree := make([]int, n+1)
   removed := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       degree[i] = len(adj[i])
   }
   q := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if degree[i] == 1 {
           q = append(q, i)
       }
   }
   for idx := 0; idx < len(q); idx++ {
       u := q[idx]
       removed[u] = true
       for _, v := range adj[u] {
           if removed[v] {
               continue
           }
           degree[v]--
           if degree[v] == 1 {
               q = append(q, v)
           }
       }
   }
   isCycle := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       if !removed[i] {
           isCycle[i] = true
       }
   }
   // extract cycle nodes in order
   var cycle []int
   var start int
   for i := 1; i <= n; i++ {
       if isCycle[i] {
           start = i
           break
       }
   }
   prev := -1
   u := start
   for {
       cycle = append(cycle, u)
       var next int
       for _, v := range adj[u] {
           if !isCycle[v] || v == prev {
               continue
           }
           next = v
           break
       }
       if next == start || next == 0 {
           break
       }
       prev, u = u, next
   }
   C := len(cycle)
   idxCycle := make([]int, n+1)
   for i, v := range cycle {
       idxCycle[v] = i
   }
   // build treeAdj (exclude cycle edges)
   treeAdj = make([][]int, n+1)
   for _, e := range edges {
       u, v := e[0], e[1]
       if isCycle[u] && isCycle[v] {
           iu := idxCycle[u]
           iv := idxCycle[v]
           if (iu+1)%C == iv || (iv+1)%C == iu {
               continue
           }
       }
       treeAdj[u] = append(treeAdj[u], v)
       treeAdj[v] = append(treeAdj[v], u)
   }
   // init HLD
   parentArr = make([]int, n+1)
   depthArr = make([]int, n+1)
   heavy = make([]int, n+1)
   sizeArr = make([]int, n+1)
   head = make([]int, n+1)
   pos = make([]int, n+1)
   for i := 1; i <= n; i++ {
       heavy[i] = -1
   }
   curPos = 0
   // for each cycle node as root
   for _, r := range cycle {
       parentArr[r] = 0
       depthArr[r] = 0
       dfsSz(r)
       decompose(r, r)
   }
   // build segtrees
   treeSeg := NewSegTree(n)
   cycleSeg := NewSegTree(C)
   // tracking
   E := 0
   onCycleCnt := 0
   // process queries
   for i := 0; i < m; i++ {
       var v, u int
       fmt.Fscan(in, &v, &u)
       totalOn := 0
       arcLen := 0
       if idxCycle[v] == 0 && !isCycle[v] {
           // ensure vRoot from parent chain
       }
       // find roots
       var rv, ru int
       if isCycle[v] {
           rv = v
       } else {
           // climb to cycle
           rv = v
           for !isCycle[rv] {
               rv = parentArr[rv]
           }
       }
       if isCycle[u] {
           ru = u
       } else {
           ru = u
           for !isCycle[ru] {
               ru = parentArr[ru]
           }
       }
       if rv == ru {
           // same tree
           on := pathOp(v, u, treeSeg, 0)
           pathOp(v, u, treeSeg, 1)
           pathLen := depthArr[v] + depthArr[u] - 2*depthArr[func() int { // compute lca
               a, b := v, u
               for head[a] != head[b] {
                   if depthArr[head[a]] > depthArr[head[b]] {
                       a = parentArr[head[a]]
                   } else {
                       b = parentArr[head[b]]
                   }
               }
               if depthArr[a] < depthArr[b] {
                   return a
               }
               return b
           }()]
           E += pathLen - 2*on
       } else {
           // v->rv
           on1 := pathOp(v, rv, treeSeg, 0)
           pathOp(v, rv, treeSeg, 1)
           len1 := depthArr[v]
           totalOn += on1
           // u->ru
           on2 := pathOp(u, ru, treeSeg, 0)
           pathOp(u, ru, treeSeg, 1)
           len2 := depthArr[u]
           totalOn += on2
           // cycle arc
           ia := idxCycle[rv]
           ib := idxCycle[ru]
           cw := (ib - ia + C) % C
           ccw := (ia - ib + C) % C
           dirCW := false
           if cw < ccw {
               dirCW = true
           } else if cw == ccw {
               // tie lex
               nxtCW := cycle[(ia+1)%C]
               nxtCCW := cycle[(ia-1+C)%C]
               if nxtCW < nxtCCW {
                   dirCW = true
               }
           }
           if dirCW {
               arcLen = cw
               if arcLen > 0 {
                   l, r := ia, (ib-1+C)%C
                   if l <= r {
                       on3 := cycleSeg.query(1, 0, C-1, l, r)
                       cycleSeg.update(1, 0, C-1, l, r)
                       totalOn += on3
                       onCycleCnt += arcLen - 2*on3
                   } else {
                       on3 := cycleSeg.query(1, 0, C-1, l, C-1) + cycleSeg.query(1, 0, C-1, 0, r)
                       cycleSeg.update(1, 0, C-1, l, C-1)
                       cycleSeg.update(1, 0, C-1, 0, r)
                       totalOn += on3
                       onCycleCnt += arcLen - 2*on3
                   }
               }
           } else {
               arcLen = ccw
               if arcLen > 0 {
                   l, r := ib, (ia-1+C)%C
                   if l <= r {
                       on3 := cycleSeg.query(1, 0, C-1, l, r)
                       cycleSeg.update(1, 0, C-1, l, r)
                       totalOn += on3
                       onCycleCnt += arcLen - 2*on3
                   } else {
                       on3 := cycleSeg.query(1, 0, C-1, l, C-1) + cycleSeg.query(1, 0, C-1, 0, r)
                       cycleSeg.update(1, 0, C-1, l, C-1)
                       cycleSeg.update(1, 0, C-1, 0, r)
                       totalOn += on3
                       onCycleCnt += arcLen - 2*on3
                   }
               }
           }
           E += (len1 + len2 + arcLen) - 2*totalOn
       }
       // compute components
       comps := n - E
       if onCycleCnt == C {
           comps++
       }
       fmt.Fprintln(out, comps)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000007

var (
   n, m int
   adj [][]int
   parent, depth, heavy, head, pos, sz []int
   curPos int
   posToNode []int
   attackTime []int
   // segment tree
   st [][]int
)

func dfs(u, p int) {
   parent[u] = p
   sz[u] = 1
   maxSz := 0
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       depth[v] = depth[u] + 1
       dfs(v, u)
       if sz[v] > maxSz {
           maxSz = sz[v]
           heavy[u] = v
       }
       sz[u] += sz[v]
   }
}

func decompose(u, h int) {
   head[u] = h
   pos[u] = curPos
   posToNode[curPos] = u
   curPos++
   if heavy[u] != -1 {
       decompose(heavy[u], h)
   }
   for _, v := range adj[u] {
       if v != parent[u] && v != heavy[u] {
           decompose(v, v)
       }
   }
}

func buildST(arr []int) {
   size := 1
   for size < n*2 {
       size <<= 1
   }
   st = make([][]int, size)
   var build func(node, l, r int)
   build = func(node, l, r int) {
       if l == r {
           if l < len(arr) {
               st[node] = []int{arr[l]}
           }
           return
       }
       mid := (l + r) >> 1
       lc := node << 1
       rc := lc | 1
       build(lc, l, mid)
       build(rc, mid+1, r)
       a, b := st[lc], st[rc]
       st[node] = make([]int, len(a)+len(b))
       i, j := 0, 0
       for k := 0; k < len(st[node]); k++ {
           if j >= len(b) || (i < len(a) && a[i] <= b[j]) {
               st[node][k] = a[i]
               i++
           } else {
               st[node][k] = b[j]
               j++
           }
       }
   }
   // build on [0, n-1]
   // find tree size levels
   // build with range [0, n-1]
   build(1, 0, n-1)
}

// count <= x in [ql,qr]
func queryLE(node, l, r, ql, qr, x int) int {
   if qr < l || r < ql {
       return 0
   }
   if ql <= l && r <= qr {
       // count from st[node]
       arr := st[node]
       // upper bound x
       lo, hi := 0, len(arr)
       for lo < hi {
           mid := (lo + hi) >> 1
           if arr[mid] <= x {
               lo = mid + 1
           } else {
               hi = mid
           }
       }
       return lo
   }
   mid := (l + r) >> 1
   return queryLE(node<<1, l, mid, ql, qr, x) + queryLE(node<<1|1, mid+1, r, ql, qr, x)
}

// count bad in segment [l,r] for y,t
func countBad(l, r, y, t int) int {
   if l > r {
       return 0
   }
   // bad: attackTime in [y+1, t]
   c1 := queryLE(1, 0, n-1, l, r, t)
   c2 := queryLE(1, 0, n-1, l, r, y)
   return c1 - c2
}

// LCA via parent and head, pos
func lca(u, v int) int {
   for head[u] != head[v] {
       if depth[head[u]] > depth[head[v]] {
           u = parent[head[u]]
       } else {
           v = parent[head[v]]
       }
   }
   if depth[u] < depth[v] {
       return u
   }
   return v
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]int, n)
   parent = make([]int, n)
   depth = make([]int, n)
   heavy = make([]int, n)
   head = make([]int, n)
   pos = make([]int, n)
   sz = make([]int, n)
   posToNode = make([]int, n)
   for i := 0; i < n; i++ {
       heavy[i] = -1
   }
   var p, root int
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p)
       if p > 0 {
           adj[i] = append(adj[i], p-1)
           adj[p-1] = append(adj[p-1], i)
       } else {
           root = i
       }
   }
   dfs(root, -1)
   curPos = 0
   decompose(root, root)
   // read events
   fmt.Fscan(in, &m)
   attackTime = make([]int, n)
   for i := 0; i < n; i++ {
       attackTime[i] = INF
   }
   type Query struct{ a, b, k, y, t, idx int }
   var queries []Query
   answers := make([]int, 0, m)
   for year := 1; year <= m; year++ {
       var tp int
       fmt.Fscan(in, &tp)
       if tp == 1 {
           var c int
           fmt.Fscan(in, &c)
           attackTime[c-1] = year
       } else {
           var a, b, k, y int
           fmt.Fscan(in, &a, &b, &k, &y)
           queries = append(queries, Query{a - 1, b - 1, k, y, year, len(answers)})
           answers = append(answers, -1)
       }
   }
   // build base array
   base := make([]int, n)
   for u := 0; u < n; u++ {
       base[pos[u]] = attackTime[u]
   }
   buildST(base)
   // process queries
   for _, q := range queries {
       a, b, k, y, t := q.a, q.b, q.k, q.y, q.t
       // build segments
       l := lca(a, b)
       type Seg struct{ l, r int; rev bool }
       var segs []Seg
       // up from a to l (exclude a)
       u := a
       for head[u] != head[l] {
           segs = append(segs, Seg{pos[head[u]], pos[u], true})
           u = parent[head[u]]
       }
       if u != l {
           segs = append(segs, Seg{pos[l] + 1, pos[u], true})
       }
       // down from l to b excluding l and b
       var down []Seg
       v := b
       for head[v] != head[l] {
           down = append(down, Seg{pos[head[v]], pos[v], false})
           v = parent[head[v]]
       }
       if v != l {
           down = append(down, Seg{pos[l] + 1, pos[v], false})
       }
       // reverse down
       for i := len(down) - 1; i >= 0; i-- {
           segs = append(segs, down[i])
       }
       // total bad and length
       totBad := 0
       totLen := 0
       for _, s := range segs {
           length := s.r - s.l + 1
           if length <= 0 {
               continue
           }
           totLen += length
           totBad += countBad(s.l, s.r, y, t)
       }
       totGood := totLen - totBad
       if k > totGood {
           answers[q.idx] = -1
           continue
       }
       // find segment
       rem := k
       found := -1
       for _, s := range segs {
           length := s.r - s.l + 1
           if length <= 0 {
               continue
           }
           bad := countBad(s.l, s.r, y, t)
           good := length - bad
           if rem > good {
               rem -= good
               continue
           }
           // binary search in this segment
           l0, r0 := s.l, s.r
           low, high := 0, r0-l0
           for low < high {
               mid := (low + high) >> 1
               var sl, sr int
               if s.rev {
                   sl = r0 - mid
                   sr = r0
               } else {
                   sl = l0
                   sr = l0 + mid
               }
               badm := countBad(sl, sr, y, t)
               goodm := (mid + 1) - badm
               if goodm >= rem {
                   high = mid
               } else {
                   low = mid + 1
                   rem -= goodm
               }
           }
           // result pos
           var pidx int
           if s.rev {
               pidx = s.r - low
           } else {
               pidx = s.l + low
           }
           found = posToNode[pidx]
           break
       }
       if found < 0 {
           answers[q.idx] = -1
       } else {
           answers[q.idx] = found + 1
       }
   }
   // output answers
   for i, v := range answers {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}

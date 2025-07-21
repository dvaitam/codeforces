package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct { to, id int }

var (
   n, m int
   adj    [][]Edge
   parent [][]int
   depth, sz, heavy, head, pos []int
   curPos int
   bit    *BIT
   // map edge id to its child node (deeper)
   childOfEdge []int
   curVal []int
   reader *bufio.Reader
   writer *bufio.Writer
)

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   adj = make([][]Edge, n+1)
   childOfEdge = make([]int, n)
   // store edges
   u := make([]int, n)
   v := make([]int, n)
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &u[i], &v[i])
       adj[u[i]] = append(adj[u[i]], Edge{v[i], i})
       adj[v[i]] = append(adj[v[i]], Edge{u[i], i})
   }
   // init arrays
   parent = make([][]int, 18)
   for i := range parent {
       parent[i] = make([]int, n+1)
   }
   depth = make([]int, n+1)
   sz = make([]int, n+1)
   heavy = make([]int, n+1)
   head = make([]int, n+1)
   pos = make([]int, n+1)
   // dfs1 to compute parent, depth, sz, heavy
   dfs1(1, 0)
   // decompose
   curPos = 1
   dfs2(1, 1)
   // prepare LCA
   for k := 1; k < len(parent); k++ {
       for v := 1; v <= n; v++ {
           parent[k][v] = parent[k-1][ parent[k-1][v] ]
       }
   }
   // map edges to child nodes
   for i := 1; i < n; i++ {
       a, b := u[i], v[i]
       if parent[0][a] == b {
           childOfEdge[i] = a
       } else {
           childOfEdge[i] = b
       }
   }
   // BIT init
   bit = NewBIT(n)
   curVal = make([]int, n+1)
   // initial all black =1, except root has no edge
   for i := 2; i <= n; i++ {
       curVal[i] = 1
       bit.Add(pos[i], 1)
   }
   curVal[1] = 0
   // queries
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 || t == 2 {
           var id int
           fmt.Fscan(reader, &id)
           child := childOfEdge[id]
           newVal := 0
           if t == 1 {
               newVal = 1
           }
           if newVal != curVal[child] {
               bit.Add(pos[child], newVal - curVal[child])
               curVal[child] = newVal
           }
       } else if t == 3 {
           var a, b int
           fmt.Fscan(reader, &a, &b)
           l := lca(a, b)
           dist := depth[a] + depth[b] - 2*depth[l]
           sum := pathSum(a, b)
           if sum == dist {
               fmt.Fprintln(writer, dist)
           } else {
               fmt.Fprintln(writer, -1)
           }
       }
   }
}

func dfs1(v, p int) {
   parent[0][v] = p
   depth[v] = depth[p] + 1
   sz[v] = 1
   maxSz := 0
   for _, e := range adj[v] {
       if e.to == p {
           continue
       }
       dfs1(e.to, v)
       if sz[e.to] > maxSz {
           maxSz = sz[e.to]
           heavy[v] = e.to
       }
       sz[v] += sz[e.to]
   }
}

func dfs2(v, h int) {
   head[v] = h
   pos[v] = curPos
   curPos++
   if heavy[v] != 0 {
       dfs2(heavy[v], h)
   }
   for _, e := range adj[v] {
       if e.to != parent[0][v] && e.to != heavy[v] {
           dfs2(e.to, e.to)
       }
   }
}

func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   // lift u
   diff := depth[u] - depth[v]
   for k := 0; diff > 0; k++ {
       if diff&1 == 1 {
           u = parent[k][u]
       }
       diff >>= 1
   }
   if u == v {
       return u
   }
   for k := len(parent) - 1; k >= 0; k-- {
       pu := parent[k][u]
       pv := parent[k][v]
       if pu != pv {
           u = pu
           v = pv
       }
   }
   return parent[0][u]
}

func pathSum(u, v int) int {
   res := 0
   for head[u] != head[v] {
       if depth[ head[u] ] > depth[ head[v] ] {
           u, v = v, u
       }
       // head[v] deeper
       res += bit.Query(pos[ head[v] ], pos[v])
       v = parent[0][ head[v] ]
   }
   // same head
   if depth[u] > depth[v] {
       u, v = v, u
   }
   // sum between u+1 .. v
   if pos[u]+1 <= pos[v] {
       res += bit.Query(pos[u]+1, pos[v])
   }
   return res
}

// BIT for sum over [1..n]
type BIT struct {
   n int
   t []int
}

func NewBIT(n int) *BIT {
   return &BIT{n, make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
   for x := i; x <= b.n; x += x & -x {
       b.t[x] += v
   }
}

// Query sum over [l..r]
func (b *BIT) Query(l, r int) int {
   return b.sum(r) - b.sum(l-1)
}

func (b *BIT) sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += b.t[x]
   }
   return s
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU on tree to answer for each query (v,k): number of colors with freq >= k in subtree of v
var (
   n, m int
   colors []int
   adj [][]int
   // euler tour mapping
   flat   []int
   st, en []int
   timer  int
   sz      []int
   heavy   []int
   // DSU data
   freqColor []int
   bit       *Fenwick
   // queries
   queries [][]Query
   ans      []int
)

// Query holds k and original index
type Query struct { k, idx int }

// Fenwick tree for frequencies count
type Fenwick struct {
   n   int
   bit []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, bit: make([]int, n+1)}
}

// add value v at position i
func (f *Fenwick) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.bit[i] += v
   }
}

// sum of [1..i]
func (f *Fenwick) Sum(i int) int {
   s := 0
   if i > f.n {
       i = f.n
   }
   for ; i > 0; i -= i & -i {
       s += f.bit[i]
   }
   return s
}

// sum of [l..r]
func (f *Fenwick) Query(l, r int) int {
   if l > r {
       return 0
   }
   if l <= 1 {
       return f.Sum(r)
   }
   return f.Sum(r) - f.Sum(l-1)
}

func dfs1(u, p int) {
   timer++
   st[u] = timer
   flat[timer] = u
   sz[u] = 1
   heavy[u] = -1
   maxSz := 0
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dfs1(v, u)
       if sz[v] > maxSz {
           maxSz = sz[v]
           heavy[u] = v
       }
       sz[u] += sz[v]
   }
   en[u] = timer
}

// add node u's color
func addNode(u int) {
   c := colors[u]
   old := freqColor[c]
   if old > 0 {
       bit.Add(old, -1)
   }
   freqColor[c] = old + 1
   bit.Add(old+1, 1)
}

// remove node u's color (for resetting)
func removeNode(u int) {
   c := colors[u]
   f := freqColor[c]
   bit.Add(f, -1)
   freqColor[c] = 0
}

func dfs2(u, p int, keep bool) {
   // process light children
   for _, v := range adj[u] {
       if v == p || v == heavy[u] {
           continue
       }
       dfs2(v, u, false)
   }
   // process heavy child
   if heavy[u] != -1 {
       dfs2(heavy[u], u, true)
   }
   // merge light children's data
   for _, v := range adj[u] {
       if v == p || v == heavy[u] {
           continue
       }
       for i := st[v]; i <= en[v]; i++ {
           addNode(flat[i])
       }
   }
   // add u
   addNode(u)
   // answer queries
   for _, q := range queries[u] {
       ans[q.idx] = bit.Query(q.k, n)
   }
   if !keep {
       // cleanup subtree u
       for i := st[u]; i <= en[u]; i++ {
           removeNode(flat[i])
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   colors = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &colors[i])
   }
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u, v := 0, 0
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   queries = make([][]Query, n+1)
   ans = make([]int, m)
   for i := 0; i < m; i++ {
       v, k := 0, 0
       fmt.Fscan(reader, &v, &k)
       queries[v] = append(queries[v], Query{k: k, idx: i})
   }
   st = make([]int, n+1)
   en = make([]int, n+1)
   flat = make([]int, n+1)
   sz = make([]int, n+1)
   heavy = make([]int, n+1)
   timer = 0
   dfs1(1, 0)
   freqColor = make([]int, 100005)
   bit = NewFenwick(n)
   dfs2(1, 0, true)
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}

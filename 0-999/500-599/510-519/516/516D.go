package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct { to int; w int64 }

var (
   n int
   adj [][]edge
   dpDown, dpUp []int64
   F []int64
   parent []int
   children [][]int
   tin, tout []int
   timer int
)

func dfsDown(u, p int) {
   dpDown[u] = 0
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       parent[v] = u
       dfsDown(v, u)
       val := dpDown[v] + e.w
       if val > dpDown[u] {
           dpDown[u] = val
       }
   }
}

func dfsUp(u, p int) {
   // find top two max dpDown+weight among children
   var m1, m2 int64
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       val := dpDown[v] + e.w
       if val >= m1 {
           m2 = m1
           m1 = val
       } else if val > m2 {
           m2 = val
       }
   }
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       use := m1
       if dpDown[v]+e.w == m1 {
           use = m2
       }
       // dpUp[v] is max of dpUp[u] and use, plus edge
       dpUp[v] = use
       if dpUp[u] > dpUp[v] {
           dpUp[v] = dpUp[u]
       }
       dpUp[v] += e.w
       dfsUp(v, u)
   }
}

func buildTree(root int) {
   parent = make([]int, n)
   children = make([][]int, n)
   // BFS/DFS to set parent and children
   stack := []int{root}
   parent[root] = -1
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, e := range adj[u] {
           v := e.to
           if v == parent[u] {
               continue
           }
           parent[v] = u
           children[u] = append(children[u], v)
           stack = append(stack, v)
       }
   }
}

func dfsTin(u int) {
   tin[u] = timer
   timer++
   for _, v := range children[u] {
       dfsTin(v)
   }
   tout[u] = timer - 1
}

// BIT for range sum on Euler tour
type BIT struct{ n int; t []int }
func newBIT(n int) *BIT { return &BIT{n, make([]int, n+1)} }
func (b *BIT) update(i, v int) {
   i++ // to 1-based
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}
func (b *BIT) query(i int) int {
   // sum [0..i]
   i++
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}
func (b *BIT) rangeQuery(l, r int) int {
   if r < l {
       return 0
   }
   return b.query(r) - b.query(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]edge, n)
   for i := 0; i < n-1; i++ {
       var x, y int
       var v int64
       fmt.Fscan(in, &x, &y, &v)
       x--;
       y--;
       adj[x] = append(adj[x], edge{y, v})
       adj[y] = append(adj[y], edge{x, v})
   }
   dpDown = make([]int64, n)
   dpUp = make([]int64, n)
   parent = make([]int, n)
   // compute dpDown and dpUp from arbitrary root 0
   dfsDown(0, -1)
   dpUp[0] = 0
   dfsUp(0, -1)
   F = make([]int64, n)
   for i := 0; i < n; i++ {
       F[i] = dpDown[i]
       if dpUp[i] > F[i] {
           F[i] = dpUp[i]
       }
   }
   // find root with minimal F
   root := 0
   for i := 1; i < n; i++ {
       if F[i] < F[root] {
           root = i
       }
   }
   buildTree(root)
   tin = make([]int, n)
   tout = make([]int, n)
   timer = 0
   dfsTin(root)
   // prepare sorted nodes by F
   type nodeF struct{ f int64; u int }
   nodes := make([]nodeF, n)
   for i := 0; i < n; i++ {
       nodes[i] = nodeF{F[i], i}
   }
   sort.Slice(nodes, func(i, j int) bool { return nodes[i].f < nodes[j].f })
   // map for order
   order := make([]int, n)
   for i, nf := range nodes {
       order[i] = nf.u
   }
   // process queries
   var q int
   fmt.Fscan(in, &q)
   ls := make([]int64, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(in, &ls[i])
   }
   // for each l
   for _, l := range ls {
       bit := newBIT(n)
       ans := 0
       j := 0
       for i := 0; i < n; i++ {
           u := order[i]
           thr := F[u] + l
           for j < n && nodes[j].f <= thr {
               bit.update(tin[nodes[j].u], 1)
               j++
           }
           // count in subtree of u
           cnt := bit.rangeQuery(tin[u], tout[u])
           if cnt > ans {
               ans = cnt
           }
       }
       fmt.Fprintln(out, ans)
   }
}

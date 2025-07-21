package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   adj    [][]int
   parent []int
   depth  []int
   dp1    []int
   dpUp   []int
   t1     []int
   max1v  []int
   max2v  []int
   par    [][]int
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent = make([]int, n+1)
   depth = make([]int, n+1)
   dp1 = make([]int, n+1)
   dpUp = make([]int, n+1)
   t1 = make([]int, n+1)
   max1v = make([]int, n+1)
   max2v = make([]int, n+1)
   // BFS for parent and depth and order
   order := make([]int, 0, n)
   q := make([]int, n)
   head, tail := 0, 0
   q[tail] = 1; tail++
   parent[1] = 0
   depth[1] = 0
   for head < tail {
       u := q[head]; head++
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           q[tail] = v; tail++
       }
   }
   // dp1: max depth in subtree
   for i := n-1; i >= 0; i-- {
       u := order[i]
       dp1[u] = 0
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           d := dp1[v] + 1
           if d > dp1[u] {
               dp1[u] = d
           }
       }
   }
   // dpUp and t1, max1v, max2v
   dpUp[1] = 0
   for _, u := range order {
       // find top two B values among neighbors, at least zero
       best1, best2, who := 0, 0, 0
       // consider parent side
       if parent[u] != 0 {
           b := dpUp[u]
           if b >= best1 {
               best2 = best1
               best1 = b
               who = parent[u]
           } else if b > best2 {
               best2 = b
           }
       }
       // consider children sides
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           b := dp1[v] + 1
           if b >= best1 {
               best2 = best1
               best1 = b
               who = v
           } else if b > best2 {
               best2 = b
           }
       }
       max1v[u] = best1
       max2v[u] = best2
       t1[u] = who
       // set dpUp for children
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           // v is child
           use := best1
           if v == who {
               use = best2
           }
           dpUp[v] = use + 1
       }
   }
   // LCA preprocessing
   LG := 0
   for (1<<uint(LG)) <= n {
       LG++
   }
   par = make([][]int, LG)
   par[0] = make([]int, n+1)
   for i := 1; i <= n; i++ {
       par[0][i] = parent[i]
   }
   for k := 1; k < LG; k++ {
       par[k] = make([]int, n+1)
       for i := 1; i <= n; i++ {
           par[k][i] = par[k-1][ par[k-1][i] ]
       }
   }
   // process queries
   var m int
   fmt.Fscan(in, &m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       ans := solveQuery(u, v)
       fmt.Fprintln(out, ans)
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
           u = par[k][u]
       }
       diff >>= 1
   }
   if u == v {
       return u
   }
   for k := len(par)-1; k >= 0; k-- {
       pu, pv := par[k][u], par[k][v]
       if pu != pv {
           u = pu
           v = pv
       }
   }
   return parent[u]
}

// ancestor returns the k-th ancestor of u (k steps up)
func ancestor(u, k int) int {
   for i := 0; k > 0; i++ {
       if k&1 == 1 {
           u = par[i][u]
       }
       k >>= 1
   }
   return u
}

// getKthNode: k-th node on path from u to v (0-indexed, 0 => u)
func getKthNode(u, v, k int) int {
   w := lca(u, v)
   du := depth[u] - depth[w]
   dv := depth[v] - depth[w]
   if k <= du {
       return ancestor(u, k)
   }
   // go down from lca towards v: steps = k-du
   down := k - du
   // remaining steps up from v = dv - down
   return ancestor(v, dv-down)
}

func solveQuery(u, v int) int {
   // distance
   w := lca(u, v)
   L := depth[u] + depth[v] - 2*depth[w]
   k := L / 2
   // wnode and tnode
   wnode := getKthNode(u, v, k)
   tnode := getKthNode(u, v, k+1)
   Fu := k + getH(wnode, tnode)
   Fv := (L - k - 1) + getH(tnode, wnode)
   if Fu > Fv {
       return Fu
   }
   return Fv
}

// getH returns h[u][t]: max distance from u excluding branch to t
func getH(u, t int) int {
   if t == t1[u] {
       return max2v[u]
   }
   return max1v[u]
}

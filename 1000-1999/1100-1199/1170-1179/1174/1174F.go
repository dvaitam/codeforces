package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)
// helper to read int
func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}
// flush output
func fl() { writer.Flush() }
// ask distance query
func askD(u int) int {
   fmt.Fprintf(writer, "? d %d\n", u)
   fl()
   return readInt()
}
// ask step query
func askS(u int) int {
   fmt.Fprintf(writer, "? s %d\n", u)
   fl()
   return readInt()
}

func main() {
   defer fl()
   n := readInt()
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt(); v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // preprocess LCA
   LOG := 18
   parent := make([][]int, LOG)
   for i := 0; i < LOG; i++ {
       parent[i] = make([]int, n+1)
   }
   depth := make([]int, n+1)
   // DFS stack
   st := []int{1}
   parent[0][1] = 0
   depth[1] = 0
   // mark visited
   vis := make([]bool, n+1)
   vis[1] = true
   for len(st) > 0 {
       u := st[len(st)-1]; st = st[:len(st)-1]
       for _, v := range adj[u] {
           if !vis[v] {
               vis[v] = true
               parent[0][v] = u
               depth[v] = depth[u] + 1
               st = append(st, v)
           }
       }
   }
   for k := 1; k < LOG; k++ {
       for v := 1; v <= n; v++ {
           parent[k][v] = parent[k-1][ parent[k-1][v] ]
       }
   }
   // lca
   lca := func(u, v int) int {
       if depth[u] < depth[v] { u, v = v, u }
       d := depth[u] - depth[v]
       for k := 0; k < LOG; k++ {
           if (d>>k)&1 != 0 {
               u = parent[k][u]
           }
       }
       if u == v { return u }
       for k := LOG-1; k >= 0; k-- {
           if parent[k][u] != parent[k][v] {
               u = parent[k][u]
               v = parent[k][v]
           }
       }
       return parent[0][u]
   }
   // distance
   dist := func(u, v int) int {
       w := lca(u, v)
       return depth[u] + depth[v] - 2*depth[w]
   }
   // k-th node on path u->v (0-based at u)
   kth := func(u, v, k int) int {
       w := lca(u, v)
       duw := depth[u] - depth[w]
       if k <= duw {
           // go up k
           x := u
           for b := 0; b < LOG; b++ {
               if (k>>b)&1 != 0 {
                   x = parent[b][x]
               }
           }
           return x
       }
       // go down from w to v
       k2 := depth[v] - depth[w] - (k - duw)
       x := v
       for b := 0; b < LOG; b++ {
           if (k2>>b)&1 != 0 {
               x = parent[b][x]
           }
       }
       return x
   }
   // active set
   active := make([]bool, n+1)
   for i := 1; i <= n; i++ { active[i] = true }
   // current u and D
   u := 1
   D := askD(1)
   // iterative search
   for {
       // compute sizes and centroid
       sz := make([]int, n+1)
       // compute subtree sizes of active from u
       var dfsSz func(int, int)
       dfsSz = func(v, p int) {
           sz[v] = 1
           for _, w := range adj[v] {
               if w == p || !active[w] { continue }
               dfsSz(w, v)
               sz[v] += sz[w]
           }
       }
       dfsSz(u, 0)
       total := sz[u]
       // find centroid
       cen := u
       changed := true
       for changed {
           changed = false
           for _, w := range adj[cen] {
               if active[w] && sz[w] < sz[cen] && sz[w] > total/2 {
                   cen = w
                   changed = true
                   break
               }
           }
       }
       // query dist at centroid
       dc := askD(cen)
       if dc == 0 {
           fmt.Fprintf(writer, "! %d\n", cen)
           return
       }
       // compute overlap t
       a := dist(u, cen)
       t := (a + D - dc) / 2
       // find w on path
       w := kth(u, cen, t)
       // query s on w
       v1 := askS(w)
       // update D and u
       D -= t + 1
       u = v1
       if D == 0 {
           fmt.Fprintf(writer, "! %d\n", u)
           return
       }
       // prune active to reachable from u
       newAct := make([]bool, n+1)
       var dfsMark func(int, int)
       dfsMark = func(v, p int) {
           newAct[v] = true
           for _, w2 := range adj[v] {
               if w2 == p || !active[w2] { continue }
               dfsMark(w2, v)
           }
       }
       dfsMark(u, 0)
       active = newAct
   }
}

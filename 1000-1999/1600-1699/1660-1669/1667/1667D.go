package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair of vertices representing an edge
type P struct{ u, v int }

var (
   g          [][]int
   depth, par, from []int
   ans, extraAns    []P
   in  = bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
)

// makeAns constructs removal moves along path between u and v with LCA a
func makeAns(u, v, a int) {
   var path []int
   cur := v
   for cur != a {
       path = append(path, cur)
       cur = par[cur]
   }
   path = append(path, a)
   var pathU []int
   cur = u
   for cur != a {
       pathU = append(pathU, cur)
       cur = par[cur]
   }
   // reverse pathU
   for i, j := 0, len(pathU)-1; i < j; i, j = i+1, j-1 {
       pathU[i], pathU[j] = pathU[j], pathU[i]
   }
   path = append(path, pathU...)
   // alternate edges starting at even index
   for i := 0; i+1 < len(path); i += 2 {
       ans = append(ans, P{path[i], path[i+1]})
   }
   // then starting at odd index
   for i := 1; i+1 < len(path); i += 2 {
       ans = append(ans, P{path[i], path[i+1]})
   }
}

// makeExtraAns constructs alternate removal moves
func makeExtraAns(u, v, a int) {
   var path []int
   cur := v
   for cur != a {
       path = append(path, cur)
       cur = par[cur]
   }
   path = append(path, a)
   var pathU []int
   cur = u
   for cur != a {
       pathU = append(pathU, cur)
       cur = par[cur]
   }
   for i, j := 0, len(pathU)-1; i < j; i, j = i+1, j-1 {
       pathU[i], pathU[j] = pathU[j], pathU[i]
   }
   path = append(path, pathU...)
   // starting at odd index
   for i := 1; i+1 < len(path); i += 2 {
       extraAns = append(extraAns, P{path[i], path[i+1]})
   }
   // then even index
   for i := 0; i+1 < len(path); i += 2 {
       extraAns = append(extraAns, P{path[i], path[i+1]})
   }
}

// dfs processes subtree rooted at v, returns false on impossible
func dfs(v, p, d int) bool {
   depth[v] = d
   par[v] = p
   vs0 := []int{}
   vs1 := []int{}
   for _, to := range g[v] {
       if to == p {
           continue
       }
       if !dfs(to, v, d+1) {
           return false
       }
       mod := (depth[from[to]] - depth[v]) & 1
       if mod == 0 {
           vs0 = append(vs0, from[to])
       } else {
           vs1 = append(vs1, from[to])
       }
   }
   // pair up vs0 and vs1
   for len(vs0) > 0 && len(vs1) > 0 {
       a := vs0[len(vs0)-1]
       vs0 = vs0[:len(vs0)-1]
       b := vs1[len(vs1)-1]
       vs1 = vs1[:len(vs1)-1]
       makeAns(a, b, v)
   }
   if len(vs0) == 0 && len(vs1) == 0 {
       from[v] = v
       return true
   }
   if len(vs0) >= 2 || len(vs1) >= 3 {
       return false
   }
   // root case
   if p < 0 {
       if len(vs0) > 0 || len(vs1) > 1 {
           return false
       }
       makeAns(v, vs1[0], v)
       return true
   }
   if len(vs0) > 0 {
       from[v] = vs0[0]
   }
   if len(vs1) > 0 {
       from[v] = vs1[0]
       if len(vs1) == 2 {
           makeExtraAns(vs1[1], v, v)
       }
   }
   return true
}

// solve handles one test case
func solve() {
   var n int
   fmt.Fscan(in, &n)
   g = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   depth = make([]int, n)
   par = make([]int, n)
   from = make([]int, n)
   for i := 0; i < n; i++ {
       from[i] = -1
   }
   ans = make([]P, 0, n-1)
   extraAns = make([]P, 0, n-1)
   ok := dfs(0, -1, 0)
   if !ok {
       fmt.Fprintln(out, "NO")
       return
   }
   fmt.Fprintln(out, "YES")
   // reverse ans
   for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
       ans[i], ans[j] = ans[j], ans[i]
   }
   // print extraAns then ans
   for _, p := range extraAns {
       fmt.Fprintln(out, p.u+1, p.v+1)
   }
   for _, p := range ans {
       fmt.Fprintln(out, p.u+1, p.v+1)
   }
}

func main() {
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       solve()
   }
}

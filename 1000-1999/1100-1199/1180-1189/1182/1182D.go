package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   if n == 1 {
       fmt.Fprintln(out, 1)
       return
   }
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--;
       v--;
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   maxDeg := 0
   for i := 0; i < n; i++ {
       if d := len(adj[i]); d > maxDeg {
           maxDeg = d
       }
   }
   if maxDeg <= 2 {
       for i := 0; i < n; i++ {
           if len(adj[i]) == 1 {
               fmt.Fprintln(out, i+1)
               return
           }
       }
   }
   // find diameter endpoints
   u := bfsFurthest(0, adj)
   parent := make([]int, n)
   for i := range parent {
       parent[i] = -1
   }
   v := bfsWithParent(u, adj, parent)
   // build diameter path
   var dia []int
   for x := v; x != -1; x = parent[x] {
       dia = append(dia, x)
   }
   // dia is from v to u, reverse to from u to v
   for i, j := 0, len(dia)-1; i < j; i, j = i+1, j-1 {
       dia[i], dia[j] = dia[j], dia[i]
   }
   // candidates: endpoints and middle
   candidates := []int{dia[0], dia[len(dia)-1], dia[len(dia)/2]}
   for _, c := range candidates {
       if solve(c, adj, n) {
           fmt.Fprintln(out, c+1)
           return
       }
   }
   // collect leaves off diameter
   mid := dia[len(dia)/2]
   isDia := make([]bool, n)
   for _, x := range dia {
       if x != mid {
           isDia[x] = true
       }
   }
   cl := collectLeaves(mid, adj, isDia)
   for _, x := range cl {
       if solve(x, adj, n) {
           fmt.Fprintln(out, x+1)
           return
       }
   }
   fmt.Fprintln(out, -1)
}

// bfsFurthest returns furthest node from start
func bfsFurthest(start int, adj [][]int) int {
   n := len(adj)
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   q := make([]int, 0, n)
   q = append(q, start)
   dist[start] = 0
   fur := start
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               q = append(q, v)
               fur = v
           }
       }
   }
   return fur
}

// bfsWithParent returns furthest node from start and fills parent
func bfsWithParent(start int, adj [][]int, parent []int) int {
   n := len(adj)
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   q := make([]int, 0, n)
   q = append(q, start)
   dist[start] = 0
   fur := start
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               parent[v] = u
               q = append(q, v)
               fur = v
           }
       }
   }
   return fur
}

// solve checks if rooting at r gives a perfect level-degree tree
func solve(r int, adj [][]int, n int) bool {
   dd := make([]int, n)
   for i := range dd {
       dd[i] = -1
   }
   vis := make([]bool, n)
   q := make([]int, 0, n)
   dQ := make([]int, 0, n)
   q = append(q, r)
   dQ = append(dQ, 0)
   vis[r] = true
   dd[0] = len(adj[r])
   for i := 0; i < len(q); i++ {
       u := q[i]
       d := dQ[i]
       for _, v := range adj[u] {
           if vis[v] {
               continue
           }
           vis[v] = true
           nd := d + 1
           dv := len(adj[v])
           if dd[nd] == -1 {
               dd[nd] = dv
           } else if dd[nd] != dv {
               return false
           }
           q = append(q, v)
           dQ = append(dQ, nd)
       }
   }
   return true
}

// collectLeaves returns leaves reachable from mid, excluding isDia nodes
func collectLeaves(mid int, adj [][]int, isDia []bool) []int {
   n := len(adj)
   vis := make([]bool, n)
   q := make([]int, 0, n)
   dQ := make([]int, 0, n)
   q = append(q, mid)
   dQ = append(dQ, 0)
   vis[mid] = true
   usedDepth := make(map[int]bool)
   var res []int
   for i := 0; i < len(q); i++ {
       u := q[i]
       d := dQ[i]
       if len(adj[u]) == 1 && u != mid {
           if !usedDepth[d] {
               usedDepth[d] = true
               res = append(res, u)
           }
       }
       for _, v := range adj[u] {
           if vis[v] || isDia[v] {
               continue
           }
           vis[v] = true
           q = append(q, v)
           dQ = append(dQ, d+1)
       }
   }
   return res
}

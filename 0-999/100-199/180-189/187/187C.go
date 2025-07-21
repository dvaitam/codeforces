package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct{ to, w int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   vols := make([]int, k)
   volIdx := make([]int, n+1)
   for i := 1; i <= n; i++ {
       volIdx[i] = -1
   }
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &vols[i])
       volIdx[vols[i]] = i
   }
   adj := make([][]int, n+1)
   var u, v int
   edgesU := make([]int, m)
   edgesV := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       edgesU[i], edgesV[i] = u, v
   }
   var s, t int
   fmt.Fscan(in, &s, &t)
   // BFS from t
   dt := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dt[i] = -1
   }
   q1 := make([]int, 0, n)
   dt[t] = 0
   q1 = append(q1, t)
   for qi := 0; qi < len(q1); qi++ {
       x := q1[qi]
       for _, y := range adj[x] {
           if dt[y] == -1 {
               dt[y] = dt[x] + 1
               q1 = append(q1, y)
           }
       }
   }
   if dt[s] < 0 {
       fmt.Println(-1)
       return
   }
   // Multi-source BFS from volunteers
   dvol := make([]int, n+1)
   leader := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dvol[i] = -1
   }
   q2 := make([]int, 0, n)
   for i, node := range vols {
       dvol[node] = 0
       leader[node] = i
       q2 = append(q2, node)
   }
   for qi := 0; qi < len(q2); qi++ {
       x := q2[qi]
       for _, y := range adj[x] {
           if dvol[y] == -1 {
               dvol[y] = dvol[x] + 1
               leader[y] = leader[x]
               q2 = append(q2, y)
           }
       }
   }
   // Build volunteers graph
   adjVol := make([][]edge, k)
   for i := 0; i < m; i++ {
       u, v = edgesU[i], edgesV[i]
       lu, lv := leader[u], leader[v]
       if lu != lv {
           w := dvol[u] + 1 + dvol[v]
           adjVol[lu] = append(adjVol[lu], edge{lv, w})
           adjVol[lv] = append(adjVol[lv], edge{lu, w})
       }
   }
   // dt for volunteers
   dtVol := make([]int, k)
   for i, node := range vols {
       dtVol[i] = dt[node]
   }
   sIdx := volIdx[s]
   // binary search on q
   lo, hi := 0, dt[s]
   for lo < hi {
       mid := (lo + hi) / 2
       if can(mid, sIdx, adjVol, dtVol) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}

func can(q, sIdx int, adjVol [][]edge, dtVol []int) bool {
   k := len(adjVol)
   vis := make([]bool, k)
   queue := make([]int, 0, k)
   vis[sIdx] = true
   queue = append(queue, sIdx)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       if dtVol[u] >= 0 && dtVol[u] <= q {
           return true
       }
       for _, e := range adjVol[u] {
           if !vis[e.to] && e.w <= q {
               vis[e.to] = true
               queue = append(queue, e.to)
           }
       }
   }
   return false
}

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

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n, a, b int
       fmt.Fscan(in, &n, &a, &b)
       adj := make([][]int, n+1)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(in, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // find path from a to b
       parent := make([]int, n+1)
       vis := make([]bool, n+1)
       queue := make([]int, 0, n)
       queue = append(queue, a)
       vis[a] = true
       parent[a] = 0
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           if u == b {
               break
           }
           for _, v := range adj[u] {
               if !vis[v] {
                   vis[v] = true
                   parent[v] = u
                   queue = append(queue, v)
               }
           }
       }
       // reconstruct path
       path := make([]int, 0, n)
       for v := b; v != 0; v = parent[v] {
           path = append(path, v)
       }
       // reverse
       L := len(path)
       for i := 0; i < L/2; i++ {
           path[i], path[L-1-i] = path[L-1-i], path[i]
       }
       // mark used nodes (path and processed branches)
       used := make([]bool, n+1)
       for _, u := range path {
           used[u] = true
       }
       // compute branch depths for each path node
       depths := make([]int, L)
       // BFS buffers
       bfsQ := make([]int, 0, n)
       bfsD := make([]int, 0, n)
       for i, u := range path {
           for _, v := range adj[u] {
               if !used[v] {
                   // BFS this branch
                   maxd := 0
                   bfsQ = bfsQ[:0]
                   bfsD = bfsD[:0]
                   bfsQ = append(bfsQ, v)
                   bfsD = append(bfsD, 0)
                   used[v] = true
                   for qi := 0; qi < len(bfsQ); qi++ {
                       x := bfsQ[qi]
                       d := bfsD[qi]
                       if d > maxd {
                           maxd = d
                       }
                       for _, w := range adj[x] {
                           if !used[w] {
                               used[w] = true
                               bfsQ = append(bfsQ, w)
                               bfsD = append(bfsD, d+1)
                           }
                       }
                   }
                   // branch depth from u includes edge u-v
                   depths[i] = max(depths[i], maxd+1)
               }
           }
       }
       // need an endpoint branch to start and an interior pivot to reverse
       if len(adj[a]) <= 1 && len(adj[b]) <= 1 {
           fmt.Fprintln(out, "NO")
           continue
       }
       ok := false
       for i := 1; i+1 < L; i++ {
           d1 := i
           d2 := L - 1 - i
           if depths[i] > d1 && depths[i] > d2 {
               ok = true
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

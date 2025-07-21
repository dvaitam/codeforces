package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([]byte, n*m)
   totalOn := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(in, &x)
           if x == 1 {
               grid[i*m+j] = 1
               totalOn++
           }
       }
   }
   if totalOn == 0 {
       fmt.Println(-1)
       return
   }
   // build adjacency and degree
   deg := make([]int, n*m)
   adj := make([][4]int, n*m)
   // for adj, we store up to 4 neighbors, fill unused with -1
   for i := range adj {
       for k := 0; k < 4; k++ {
           adj[i][k] = -1
       }
   }
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i*m+j] == 0 {
               continue
           }
           id := i*m + j
           cnt := 0
           for d := 0; d < 4; d++ {
               ni, nj := i+dirs[d][0], j+dirs[d][1]
               if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni*m+nj] == 1 {
                   adj[id][cnt] = ni*m + nj
                   cnt++
               }
           }
           deg[id] = cnt
           if cnt > 2 {
               fmt.Println(-1)
               return
           }
       }
   }
   // connectivity BFS
   vis := make([]bool, n*m)
   var q []int
   // find first one
   start := -1
   for i := 0; i < n*m; i++ {
       if grid[i] == 1 {
           start = i
           break
       }
   }
   vis[start] = true
   q = append(q, start)
   cntVis := 1
   for head := 0; head < len(q); head++ {
       u := q[head]
       for k := 0; k < deg[u]; k++ {
           v := adj[u][k]
           if !vis[v] {
               vis[v] = true
               q = append(q, v)
               cntVis++
           }
       }
   }
   if cntVis != totalOn {
       fmt.Println(-1)
       return
   }
   // count edges and endpoints
   edgeCount := 0
   endpoints := []int{}
   for i := 0; i < n*m; i++ {
       if grid[i] == 1 {
           edgeCount += deg[i]
           if deg[i] == 1 {
               endpoints = append(endpoints, i)
           }
       }
   }
   edgeCount /= 2
   if edgeCount == 0 {
       fmt.Println(-1)
       return
   }
   isCycle := false
   if len(endpoints) == 0 {
       // possible cycle
       isCycle = true
   } else if len(endpoints) == 2 {
       // path
       isCycle = false
       start = endpoints[0]
   } else {
       fmt.Println(-1)
       return
   }
   if !isCycle && len(endpoints) == 2 {
       start = endpoints[0]
   }
   // traverse edges in order
   path := make([]int, 0, edgeCount+1)
   prev := -1
   curr := start
   path = append(path, curr)
   // if cycle, we still set prev=-1 and stop after edgeCount steps
   for steps := 0; steps < edgeCount; steps++ {
       // find next neighbor not prev
       found := false
       for k := 0; k < deg[curr]; k++ {
           v := adj[curr][k]
           if v != prev {
               // next
               prev, curr = curr, v
               path = append(path, curr)
               found = true
               break
           }
       }
       if !found {
           break
       }
   }
   // compute runs
   G := 0
   // if path length is less than 2, no runs
   if len(path) < 2 {
       fmt.Println(-1)
       return
   }
   // compute direction of first edge
   prevDir := [2]int{path[1]/m - path[0]/m, path[1]%m - path[0]%m}
   runLen := 1
   for i := 1; i < len(path)-1; i++ {
       dx := path[i+1]/m - path[i]/m
       dy := path[i+1]%m - path[i]%m
       if dx == prevDir[0] && dy == prevDir[1] {
           runLen++
       } else {
           G = gcd(G, runLen)
           runLen = 1
           prevDir = [2]int{dx, dy}
       }
   }
   G = gcd(G, runLen)
   // collect divisors >1
   res := []int{}
   for k := 2; k*k <= G; k++ {
       if G%k == 0 {
           res = append(res, k)
           if k != G/k {
               res = append(res, G/k)
           }
       }
   }
   if G > 1 {
       res = append(res, G)
   }
   if len(res) == 0 {
       fmt.Println(-1)
       return
   }
   sort.Ints(res)
   // print unique
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range res {
       if i > 0 {
           out.WriteString(" ")
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}

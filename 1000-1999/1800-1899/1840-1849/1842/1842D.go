package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(1e18)

type Edge struct {
   to int
   w  int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   edges := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       var y int64
       fmt.Fscan(reader, &u, &v, &y)
       edges[u] = append(edges[u], Edge{to: v, w: y})
       edges[v] = append(edges[v], Edge{to: u, w: y})
   }
   // Dijkstra from node n
   dist := make([]int64, n+1)
   vis := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[n] = 0
   for {
       u := -1
       best := INF + 1
       for i := 1; i <= n; i++ {
           if !vis[i] && dist[i] < best {
               best = dist[i]
               u = i
           }
       }
       if u == -1 {
           break
       }
       vis[u] = true
       for _, e := range edges[u] {
           if dist[e.to] > dist[u]+e.w {
               dist[e.to] = dist[u] + e.w
           }
       }
   }
   // Initialize masks: true if dist>0 (including unreachable INF), false if dist==0
   mask := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       if dist[i] == 0 {
           mask[i] = false
       } else {
           mask[i] = true
       }
   }
   // Collect positive finite distances
   dset := make(map[int64]struct{})
   for i := 1; i <= n; i++ {
       if dist[i] > 0 && dist[i] < INF {
           dset[dist[i]] = struct{}{}
       }
   }
   dists := make([]int64, 0, len(dset))
   for d := range dset {
       dists = append(dists, d)
   }
   sort.Slice(dists, func(i, j int) bool { return dists[i] < dists[j] })
   type Res struct {
       s  string
       qt int64
   }
   var resp []Res
   var rs int64
   // Generate games
   for _, d := range dists {
       if !mask[1] {
           break
       }
       // game with current mask for duration d - rs
       sb := make([]byte, n)
       for i := 1; i <= n; i++ {
           if mask[i] {
               sb[i-1] = '1'
           } else {
               sb[i-1] = '0'
           }
       }
       qt := d - rs
       resp = append(resp, Res{s: string(sb), qt: qt})
       // turn off all with dist == d
       for i := 1; i <= n; i++ {
           if dist[i] == d {
               mask[i] = false
           }
       }
       rs = d
   }
   // Output
   if mask[1] {
       fmt.Fprintln(writer, "inf")
   } else {
       fmt.Fprintf(writer, "%d %d\n", rs, len(resp))
       for _, r := range resp {
           fmt.Fprintf(writer, "%s %d\n", r.s, r.qt)
       }
   }
}

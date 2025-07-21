package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1e9

type pair struct { c, d int }

var (
   n, m        int
   adj         [][]int
   size        []int
   used        []bool
   parentCent  []int
   dists       [][]pair
   bestDist    []int
)

func calcSize(u, p int) int {
   size[u] = 1
   for _, v := range adj[u] {
       if v != p && !used[v] {
           size[u] += calcSize(v, u)
       }
   }
   return size[u]
}

func findCentroid(u, p, tot int) int {
   for _, v := range adj[u] {
       if v != p && !used[v] && size[v] > tot/2 {
           return findCentroid(v, u, tot)
       }
   }
   return u
}

func assignDists(u, p, depth, cent int) {
   dists[u] = append(dists[u], pair{cent, depth})
   for _, v := range adj[u] {
       if v != p && !used[v] {
           assignDists(v, u, depth+1, cent)
       }
   }
}

func decompose(u, p int) {
   tot := calcSize(u, -1)
   cent := findCentroid(u, -1, tot)
   if p < 0 {
       parentCent[cent] = cent
   } else {
       parentCent[cent] = p
   }
   used[cent] = true
   assignDists(cent, -1, 0, cent)
   for _, v := range adj[cent] {
       if !used[v] {
           decompose(v, cent)
       }
   }
}

func updateRed(u int) {
   for _, pr := range dists[u] {
       if pr.d < bestDist[pr.c] {
           bestDist[pr.c] = pr.d
       }
   }
}

func queryMin(u int) int {
   res := INF
   for _, pr := range dists[u] {
       d := pr.d + bestDist[pr.c]
       if d < res {
           res = d
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   adj = make([][]int, n+1)
   size = make([]int, n+1)
   used = make([]bool, n+1)
   parentCent = make([]int, n+1)
   dists = make([][]pair, n+1)
   bestDist = make([]int, n+1)
   for i := 1; i <= n; i++ {
       bestDist[i] = INF
   }
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   decompose(1, -1)
   // initially node 1 is red
   updateRed(1)
   for i := 0; i < m; i++ {
       var t, v int
       fmt.Fscan(reader, &t, &v)
       if t == 1 {
           updateRed(v)
       } else {
           ans := queryMin(v)
           fmt.Fprintln(writer, ans)
       }
   }
}

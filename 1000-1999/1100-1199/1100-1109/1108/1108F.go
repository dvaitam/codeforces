package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct {
   u, v, w int
}

type pair struct {
   x, y int
}

var parent []int

func find(u int) int {
   if parent[u] != u {
       parent[u] = find(parent[u])
   }
   return parent[u]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   edges := make([]edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })

   parent = make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }

   cnt := 0
   prevW := -1
   pairs := make([]pair, 0)

   processGroup := func() {
       for _, p := range pairs {
           xRoot := find(p.x)
           yRoot := find(p.y)
           if xRoot != yRoot {
               parent[xRoot] = yRoot
           } else {
               cnt++
           }
       }
       pairs = pairs[:0]
   }

   for _, e := range edges {
       if e.w != prevW {
           processGroup()
           prevW = e.w
       }
       uRoot := find(e.u)
       vRoot := find(e.v)
       if uRoot != vRoot {
           if uRoot > vRoot {
               uRoot, vRoot = vRoot, uRoot
           }
           pairs = append(pairs, pair{uRoot, vRoot})
       }
   }
   processGroup()

   fmt.Fprintln(writer, cnt)
}

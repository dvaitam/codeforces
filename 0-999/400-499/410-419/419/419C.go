package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, p int
   fmt.Fscan(reader, &n, &p)
   deg := make([]int, n+1)
   type pair struct{ u, v int }
   edges := make([]pair, 0, n)

   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       deg[x]++
       deg[y]++
       u, v := x, y
       if u > v {
           u, v = v, u
       }
       edges = append(edges, pair{u, v})
   }

   // Prepare sorted degrees
   sortedDeg := make([]int, n)
   for i := 1; i <= n; i++ {
       sortedDeg[i-1] = deg[i]
   }
   sort.Ints(sortedDeg)

   // Count pairs with deg[i] + deg[j] >= p
   var ans int64
   for i := 0; i < n; i++ {
       need := p - sortedDeg[i]
       j := sort.Search(n, func(j int) bool {
           return sortedDeg[j] >= need
       })
       if j < i+1 {
           j = i + 1
       }
       if j < n {
           ans += int64(n - j)
       }
   }

   // Adjust for pairs connected by edges
   // Count common edges per pair
   sort.Slice(edges, func(i, j int) bool {
       if edges[i].u != edges[j].u {
           return edges[i].u < edges[j].u
       }
       return edges[i].v < edges[j].v
   })
   // scan edges
   for i := 0; i < len(edges); {
       j := i + 1
       for j < len(edges) && edges[j] == edges[i] {
           j++
       }
       cnt := j - i
       u, v := edges[i].u, edges[i].v
       // if this pair was counted but fails with common edges
       if deg[u]+deg[v] >= p && deg[u]+deg[v]-cnt < p {
           ans--
       }
       i = j
   }

   fmt.Fprintln(writer, ans)
}

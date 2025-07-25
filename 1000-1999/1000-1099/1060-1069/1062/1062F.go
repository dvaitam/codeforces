package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   adj := make([][]int, n+1)
   radj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       radj[v] = append(radj[v], u)
   }
   // Precompute degree sum for filtering
   degSum := make([]int, n+1)
   for u := 1; u <= n; u++ {
       degSum[u] = len(adj[u]) + len(radj[u])
   }
   var candidates []int
   for u := 1; u <= n; u++ {
       if degSum[u] >= n-2 {
           candidates = append(candidates, u)
       }
   }
   visitedF := make([]bool, n+1)
   visitedR := make([]bool, n+1)
   ans := 0
   // Evaluate each candidate by BFS
   for _, u := range candidates {
       // forward BFS to count descendants
       for i := range visitedF {
           visitedF[i] = false
       }
       var queue []int
       visitedF[u] = true
       queue = append(queue, u)
       countF := 0
       for qi := 0; qi < len(queue); qi++ {
           x := queue[qi]
           for _, v := range adj[x] {
               if !visitedF[v] {
                   visitedF[v] = true
                   countF++
                   queue = append(queue, v)
               }
           }
       }
       // backward BFS to count ancestors
       for i := range visitedR {
           visitedR[i] = false
       }
       queue = queue[:0]
       visitedR[u] = true
       queue = append(queue, u)
       countR := 0
       for qi := 0; qi < len(queue); qi++ {
           x := queue[qi]
           for _, v := range radj[x] {
               if !visitedR[v] {
                   visitedR[v] = true
                   countR++
                   queue = append(queue, v)
               }
           }
       }
       if countF+countR >= n-2 {
           ans++
       }
   }
   fmt.Println(ans)
}

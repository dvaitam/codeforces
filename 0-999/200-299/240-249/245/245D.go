package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   b := make([][]int, n)
   for i := 0; i < n; i++ {
       b[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Scan(&b[i][j])
       }
   }
   a := make([]int, n)
   // Reconstruct bits independently
   for k := 0; k <= 30; k++ {
       // degree of each node in bit-k graph
       deg := make([]int, n)
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if i != j && ((b[i][j]>>k)&1) == 1 {
                   deg[i]++
               }
           }
       }
       visited := make([]bool, n)
       for i := 0; i < n; i++ {
           if visited[i] || deg[i] == 0 {
               continue
           }
           // BFS to collect connected component of bit-k edges
           queue := []int{i}
           comp := []int{i}
           visited[i] = true
           for head := 0; head < len(queue); head++ {
               u := queue[head]
               for v := 0; v < n; v++ {
                   if !visited[v] && ((b[u][v]>>k)&1) == 1 {
                       visited[v] = true
                       queue = append(queue, v)
                       comp = append(comp, v)
                   }
               }
           }
           // Set bit k for all nodes in this clique
           for _, u := range comp {
               a[u] |= 1 << k
           }
       }
   }
   // Output result
   for i, v := range a {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}

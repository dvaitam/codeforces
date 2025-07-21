package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   visited := make([]bool, n+1)
   bench := 0
   oddPathCount := 0

   // Explore each component
   for i := 1; i <= n; i++ {
       if visited[i] {
           continue
       }
       // BFS/DFS to count nodes and sum degrees
       stack := []int{i}
       visited[i] = true
       nodes := 0
       degSum := 0
       for len(stack) > 0 {
           u := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           nodes++
           degSum += len(adj[u])
           for _, v := range adj[u] {
               if !visited[v] {
                   visited[v] = true
                   stack = append(stack, v)
               }
           }
       }
       edges := degSum / 2
       if edges == nodes {
           // cycle
           if nodes%2 == 1 {
               // odd cycle: bench one
               bench++
           }
           // even cycle contributes no imbalance
       } else {
           // tree/path component
           if nodes%2 == 1 {
               oddPathCount++
           }
       }
   }
   // If odd number of odd-sized paths, need to bench one more
   if oddPathCount%2 == 1 {
       bench++
   }
   fmt.Fprintln(writer, bench)
}

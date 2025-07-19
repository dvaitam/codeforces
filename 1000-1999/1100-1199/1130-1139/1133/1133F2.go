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
   var n, m, D int
   if _, err := fmt.Fscan(in, &n, &m, &D); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // Degree of node 1
   deg1 := len(adj[1])
   if deg1 < D {
       fmt.Fprintln(out, "NO")
       return
   }
   // Find connected components in graph without node 1
   comp := make([]int, n+1)
   compID := 0
   for i := 2; i <= n; i++ {
       if comp[i] == 0 {
           compID++
           // BFS/DFS from i, skipping node 1
           stack := []int{i}
           comp[i] = compID
           for len(stack) > 0 {
               u := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               for _, v := range adj[u] {
                   if v == 1 || comp[v] != 0 {
                       continue
                   }
                   comp[v] = compID
                   stack = append(stack, v)
               }
           }
       }
   }
   if compID > D {
       fmt.Fprintln(out, "NO")
       return
   }
   // Group neighbors of 1 by component
   compEdges := make(map[int][]int)
   for _, y := range adj[1] {
       // comp[y] must be >0
       compEdges[comp[y]] = append(compEdges[comp[y]], y)
   }
   // Select initial edges from 1
   chosen := make(map[int]bool)
   initial := make([]int, 0, D)
   // One per component
   for cid := 1; cid <= compID; cid++ {
       list := compEdges[cid]
       if len(list) == 0 {
           // No edge to this component, impossible
           fmt.Fprintln(out, "NO")
           return
       }
       y := list[0]
       chosen[y] = true
       initial = append(initial, y)
   }
   // Use extra edges to reach D
   left := D - compID
   for _, y := range adj[1] {
       if left == 0 {
           break
       }
       if chosen[y] {
           continue
       }
       chosen[y] = true
       initial = append(initial, y)
       left--
   }
   // Build spanning tree edges
   edges := make([][2]int, 0, n-1)
   // Edges from 1
   visited := make([]bool, n+1)
   visited[1] = true
   for _, y := range initial {
       edges = append(edges, [2]int{1, y})
       visited[y] = true
   }
   // BFS from initial neighbors to cover rest
   queue := make([]int, len(initial))
   copy(queue, initial)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if v == 1 || visited[v] {
               continue
           }
           visited[v] = true
           edges = append(edges, [2]int{u, v})
           queue = append(queue, v)
       }
   }
   // Output result
   fmt.Fprintln(out, "YES")
   for _, e := range edges {
       fmt.Fprintln(out, e[0], e[1])
   }
}

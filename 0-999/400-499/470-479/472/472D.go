package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to int
   w  int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   const INF = int64(1e18)
   // Read distance matrix
   d := make([][]int64, n)
   for i := 0; i < n; i++ {
       row := make([]int64, n)
       for j := 0; j < n; j++ {
           if _, err := fmt.Fscan(in, &row[j]); err != nil {
               return
           }
       }
       d[i] = row
   }
   // Basic validation
   for i := 0; i < n; i++ {
       if d[i][i] != 0 {
           fmt.Println("NO")
           return
       }
       for j := 0; j < n; j++ {
           if d[i][j] != d[j][i] {
               fmt.Println("NO")
               return
           }
           if i != j && d[i][j] == 0 {
               fmt.Println("NO")
               return
           }
       }
   }
   if n == 1 {
       fmt.Println("YES")
       return
   }
   // Prim's algorithm to build MST from complete graph
   used := make([]bool, n)
   key := make([]int64, n)
   parent := make([]int, n)
   for i := 0; i < n; i++ {
       key[i] = INF
       parent[i] = -1
   }
   key[0] = 0
   for iter := 0; iter < n; iter++ {
       u := -1
       best := INF
       for i := 0; i < n; i++ {
           if !used[i] && key[i] < best {
               best = key[i]
               u = i
           }
       }
       if u == -1 {
           fmt.Println("NO")
           return
       }
       used[u] = true
       for v := 0; v < n; v++ {
           if !used[v] && d[u][v] < key[v] {
               key[v] = d[u][v]
               parent[v] = u
           }
       }
   }
   // Build tree adjacency
   adj := make([][]edge, n)
   for v := 1; v < n; v++ {
       u := parent[v]
       if u < 0 || u >= n || key[v] <= 0 {
           fmt.Println("NO")
           return
       }
       w := key[v]
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   // Verify distances with BFS from each node
   dist := make([]int64, n)
   q := make([]int, n)
   for src := 0; src < n; src++ {
       for i := 0; i < n; i++ {
           dist[i] = -1
       }
       head, tail := 0, 0
       q[tail] = src
       tail++
       dist[src] = 0
       for head < tail {
           u := q[head]
           head++
           for _, e := range adj[u] {
               v := e.to
               if dist[v] < 0 {
                   dist[v] = dist[u] + e.w
                   q[tail] = v
                   tail++
               }
           }
       }
       for j := 0; j < n; j++ {
           if dist[j] != d[src][j] {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
}

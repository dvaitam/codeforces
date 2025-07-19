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
   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS from node 1 to find farthest a
   a, _ := bfsFarthest(1, adj)
   // BFS from a to find farthest b and record parents, distances
   b, dist, parent := bfsFarthestWithParent(a, adj)
   // mark diameter path from b to a
   vis := make([]bool, n+1)
   path := []int{}
   for x := b; x != 0; x = parent[x] {
       vis[x] = true
       path = append(path, x)
   }
   // multi-source BFS from all nodes on diameter path to find c
   queue := make([]int, 0, n)
   dist3 := make([]int, n+1)
   for _, u := range path {
       queue = append(queue, u)
       dist3[u] = 0
   }
   c, maxd := path[0], 0
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       d := dist3[u]
       if d > maxd {
           maxd = d
           c = u
       }
       for _, v := range adj[u] {
           if vis[v] {
               continue
           }
           vis[v] = true
           dist3[v] = d + 1
           queue = append(queue, v)
       }
   }
   // ensure c is distinct from a and b
   if c == a || c == b {
       for i := 1; i <= n; i++ {
           if i != a && i != b {
               c = i
               break
           }
       }
   }
   // result is diameter length + maxd
   res := dist[b] + maxd
   fmt.Fprintln(out, res)
   fmt.Fprintln(out, a, b, c)
}

// bfsFarthest returns farthest node and its distance from start
func bfsFarthest(start int, adj [][]int) (node int, dist []int) {
   n := len(adj) - 1
   dist = make([]int, n+1)
   vis := make([]bool, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, start)
   vis[start] = true
   node = start
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if vis[v] {
               continue
           }
           vis[v] = true
           dist[v] = dist[u] + 1
           queue = append(queue, v)
           if dist[v] > dist[node] {
               node = v
           }
       }
   }
   return
}

// bfsFarthestWithParent returns farthest node, dist array, and parent pointers
func bfsFarthestWithParent(start int, adj [][]int) (node int, dist, parent []int) {
   n := len(adj) - 1
   dist = make([]int, n+1)
   parent = make([]int, n+1)
   vis := make([]bool, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, start)
   vis[start] = true
   parent[start] = 0
   node = start
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if vis[v] {
               continue
           }
           vis[v] = true
           parent[v] = u
           dist[v] = dist[u] + 1
           queue = append(queue, v)
           if dist[v] > dist[node] {
               node = v
           }
       }
   }
   return
}

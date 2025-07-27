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

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n, a, b, da, db int
       fmt.Fscan(reader, &n, &a, &b, &da, &db)
       a--
       b--
       adj := make([][]int, n)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           u--
           v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // distance between a and b
       distAB := bfsDist(a, b, adj, n)
       if distAB <= da {
           fmt.Fprintln(writer, "Alice")
           continue
       }
       // compute tree diameter
       farNode, _ := bfsFarthest(0, adj, n)
       _, diameter := bfsFarthest(farNode, adj, n)
       // Alice wins if Bob cannot run far
       if db <= 2*da || diameter <= 2*da {
           fmt.Fprintln(writer, "Alice")
       } else {
           fmt.Fprintln(writer, "Bob")
       }
   }
}

// bfsDist returns distance from src to target
func bfsDist(src, target int, adj [][]int, n int) int {
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   queue := make([]int, 0, n)
   dist[src] = 0
   queue = append(queue, src)
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       if u == target {
           break
       }
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               queue = append(queue, v)
           }
       }
   }
   return dist[target]
}

// bfsFarthest returns farthest node and its distance from start
func bfsFarthest(start int, adj [][]int, n int) (node, maxd int) {
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   queue := make([]int, 0, n)
   dist[start] = 0
   queue = append(queue, start)
   node = start
   maxd = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               queue = append(queue, v)
               if dist[v] > maxd {
                   maxd = dist[v]
                   node = v
               }
           }
       }
   }
   return node, maxd
}

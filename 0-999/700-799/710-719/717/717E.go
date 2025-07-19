package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n      int
   color  []int
   a      []bool
   edges  [][]int
   writer *bufio.Writer
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   color = make([]int, n+1)
   a = make([]bool, n+1)
   edges = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &color[i])
       if color[i] == -1 {
           color[i] = 0
       }
   }
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       edges[u] = append(edges[u], v)
       edges[v] = append(edges[v], u)
   }
   dfs(1, 0)
   solve(1, 0)
   if color[1] == 0 && len(edges[1]) > 0 {
       x := edges[1][0]
       fmt.Fprintf(writer, "%d %d %d", x, 1, x)
   }
}

// dfs computes a[u] = true if u or any node in its subtree has color != 1
func dfs(u, p int) {
   if color[u] != 1 {
       a[u] = true
   }
   for _, v := range edges[u] {
       if v == p {
           continue
       }
       dfs(v, u)
       if a[v] {
           a[u] = true
       }
   }
}

// solve prints the traversal sequence according to the original logic
func solve(u, p int) {
   if u != 1 {
       color[u] ^= 1
   }
   fmt.Fprintf(writer, "%d ", u)
   for _, v := range edges[u] {
       if v == p {
           continue
       }
       if a[v] || color[v] == 0 {
           solve(v, u)
           color[u] ^= 1
           fmt.Fprintf(writer, "%d ", u)
           if color[v] == 0 {
               color[u] ^= 1
               fmt.Fprintf(writer, "%d %d ", v, u)
           }
       }
   }
}

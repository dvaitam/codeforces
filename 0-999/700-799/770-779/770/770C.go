package main

import (
   "fmt"
   "os"
)

var (
   n, k     int
   graph    [][]int
   a        []int
   vis      []int
   order    []int
   found    bool
)

func dfs(x int) {
   vis[x] = 2
   for _, y := range graph[x] {
       if found {
           return
       }
       if vis[y] == 0 {
           dfs(y)
       } else if vis[y] == 2 {
           found = true
           return
       }
   }
   vis[x] = 1
   order = append(order, x)
}

func main() {
   // Input and output via standard streams
   in := os.Stdin
   out := os.Stdout
   fmt.Fscan(in, &n, &k)
   a = make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &a[i])
   }
   graph = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       var si int
       fmt.Fscan(in, &si)
       if si > 0 {
           graph[i] = make([]int, si)
           for j := 0; j < si; j++ {
               fmt.Fscan(in, &graph[i][j])
           }
       }
   }
   vis = make([]int, n+1)
   order = make([]int, 0, n)
   for _, x := range a {
       if vis[x] == 0 {
           dfs(x)
       }
   }
   if found {
       fmt.Fprintln(out, -1)
       return
   }
   fmt.Fprintln(out, len(order))
   for _, v := range order {
       fmt.Fprint(out, v, " ")
   }
   fmt.Fprintln(out)
}

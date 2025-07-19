package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   r := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(r, &n); err != nil {
       return
   }
   // adjacency matrix for vertices 0..6
   var a [7][7]int
   x := make([]int, n)
   y := make([]int, n)
   deg := [7]int{}
   for i := 0; i < n; i++ {
       fmt.Fscan(r, &x[i], &y[i])
       a[x[i]][y[i]]++
       a[y[i]][x[i]]++
       deg[x[i]]++
       deg[y[i]]++
   }
   // find start vertex: first edge's first endpoint or a vertex with odd degree
   start := x[0]
   odd := 0
   for i := 0; i <= 6; i++ {
       if deg[i]%2 != 0 {
           odd++
           start = i
       }
   }
   if odd > 2 {
       fmt.Println("No solution")
       return
   }
   // Hierholzer's algorithm
   ans := make([]int, 0, n+1)
   var dfs func(u int)
   dfs = func(u int) {
       for v := 0; v <= 6; v++ {
           for a[u][v] > 0 {
               a[u][v]--
               a[v][u]--
               dfs(v)
           }
       }
       ans = append(ans, u)
   }
   dfs(start)
   if len(ans) != n+1 {
       fmt.Println("No solution")
       return
   }
   // reverse to get path from start to end
   for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
       ans[i], ans[j] = ans[j], ans[i]
   }
   // output edges with direction
   used := make([]bool, n)
   for i := 0; i < n; i++ {
       u, v := ans[i], ans[i+1]
       found := false
       for j := 0; j < n; j++ {
           if used[j] {
               continue
           }
           if u == x[j] && v == y[j] {
               fmt.Printf("%d +\n", j+1)
               used[j] = true
               found = true
               break
           } else if u == y[j] && v == x[j] {
               fmt.Printf("%d -\n", j+1)
               used[j] = true
               found = true
               break
           }
       }
       if !found {
           fmt.Println("No solution")
           return
       }
   }
}

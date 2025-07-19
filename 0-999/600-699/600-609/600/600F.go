package main

import (
   "bufio"
   "fmt"
   "os"
)

var nx, ny, m int
var id [][]int
var dx, dy []int
var cx, cy [][]int
var ans []int

func dfs(x [][]int, y [][]int, u, v, cu, cv int) {
   w := y[v][cv]
   if x[w][cu] == 0 {
       x[w][cu] = v
       y[v][cu] = w
       x[w][cv] = 0
   } else {
       dfs(y, x, v, w, cv, cu)
   }
   x[u][cv] = v
   y[v][cv] = u
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &nx, &ny, &m)
   id = make([][]int, nx+1)
   for i := range id {
       id[i] = make([]int, ny+1)
   }
   dx = make([]int, nx+1)
   dy = make([]int, ny+1)
   edgesU := make([]int, m+1)
   edgesV := make([]int, m+1)
   for i := 1; i <= m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       edgesU[i] = u
       edgesV[i] = v
       id[u][v] = i
       dx[u]++
       dy[v]++
   }
   d := 0
   for i := 1; i <= nx; i++ {
       if dx[i] > d {
           d = dx[i]
       }
   }
   for i := 1; i <= ny; i++ {
       if dy[i] > d {
           d = dy[i]
       }
   }
   // output number of colors
   fmt.Println(d)
   // initialize color assignment arrays
   cx = make([][]int, nx+1)
   for i := range cx {
       cx[i] = make([]int, d+1)
   }
   cy = make([][]int, ny+1)
   for i := range cy {
       cy[i] = make([]int, d+1)
   }
   // assign colors
   for u := 1; u <= nx; u++ {
       for v := 1; v <= ny; v++ {
           if id[u][v] != 0 {
               c := 1
               for c <= d && (cx[u][c] != 0 || cy[v][c] != 0) {
                   c++
               }
               if c <= d {
                   cx[u][c] = v
                   cy[v][c] = u
                   continue
               }
               cu := 1
               cv := 1
               for cu <= d && cy[v][cu] != 0 {
                   cu++
               }
               for cv <= d && cx[u][cv] != 0 {
                   cv++
               }
               dfs(cx, cy, u, v, cu, cv)
           }
       }
   }
   ans = make([]int, m+1)
   for u := 1; u <= nx; u++ {
       for c := 1; c <= d; c++ {
           v := cx[u][c]
           if v != 0 {
               idx := id[u][v]
               ans[idx] = c
           }
       }
   }
   // output colors for each edge
   w := bufio.NewWriter(os.Stdout)
   for i := 1; i <= m; i++ {
       if i > 1 {
           w.WriteByte(' ')
       }
       w.WriteString(fmt.Sprint(ans[i]))
   }
   w.WriteByte('\n')
   w.Flush()
}

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
   var n, m, q int
   if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
       return
   }
   grid := make([][]int, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
   }
   var x, y, c int
   // For BFS
   dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   visited := make([][]bool, n)
   for i := 0; i < n; i++ {
       visited[i] = make([]bool, m)
   }
   for qi := 0; qi < q; qi++ {
       fmt.Fscan(reader, &x, &y, &c)
       x--
       y--
       grid[x][y] = c
       // count components
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               visited[i][j] = false
           }
       }
       comps := 0
       var stack [100000]int
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if !visited[i][j] {
                   comps++ 
                   // DFS
                   top := 0
                   stack[top] = i*m + j
                   visited[i][j] = true
                   for top >= 0 {
                       idx := stack[top]
                       top--
                       cx := idx / m
                       cy := idx % m
                       for _, d := range dirs {
                           nx, ny := cx+d[0], cy+d[1]
                           if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] && grid[nx][ny] == grid[cx][cy] {
                               visited[nx][ny] = true
                               top++
                               stack[top] = nx*m + ny
                           }
                       }
                   }
               }
           }
       }
       fmt.Fprintln(writer, comps)
   }
}

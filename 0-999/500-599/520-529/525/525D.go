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
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   visited := make([][]bool, n)
   for i := range visited {
       visited[i] = make([]bool, m)
   }
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' && !visited[i][j] {
               // explore component
               minr, maxr := i, i
               minc, maxc := j, j
               stack := [][2]int{{i, j}}
               visited[i][j] = true
               for len(stack) > 0 {
                   r, c := stack[len(stack)-1][0], stack[len(stack)-1][1]
                   stack = stack[:len(stack)-1]
                   if r < minr {
                       minr = r
                   }
                   if r > maxr {
                       maxr = r
                   }
                   if c < minc {
                       minc = c
                   }
                   if c > maxc {
                       maxc = c
                   }
                   for _, d := range dirs {
                       nr, nc := r+d[0], c+d[1]
                       if nr >= 0 && nr < n && nc >= 0 && nc < m && !visited[nr][nc] && grid[nr][nc] == '.' {
                           visited[nr][nc] = true
                           stack = append(stack, [2]int{nr, nc})
                       }
                   }
               }
               // fill bounding rectangle
               for rr := minr; rr <= maxr; rr++ {
                   for cc := minc; cc <= maxc; cc++ {
                       grid[rr][cc] = '.'
                   }
               }
           }
       }
   }
   // output result
   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}

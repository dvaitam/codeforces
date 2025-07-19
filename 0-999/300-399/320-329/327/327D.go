package main

import (
   "bufio"
   "fmt"
   "os"
)

type Point struct { x, y int }
type Op struct { Opt byte; x, y int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([][]byte, n+2)
   for i := 0; i < n+2; i++ {
       grid[i] = make([]byte, m+2)
   }
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 1; j <= m; j++ {
           grid[i][j] = line[j-1]
       }
   }
   mark := make([][]bool, n+2)
   vis := make([][]bool, n+2)
   for i := 0; i < n+2; i++ {
       mark[i] = make([]bool, m+2)
       vis[i] = make([]bool, m+2)
   }
   var ans []Op
   ans = make([]Op, 0, n*m*2)
   var stack []Point
   stack = make([]Point, 0, n*m)
   dx := [4]int{0, 0, 1, -1}
   dy := [4]int{1, -1, 0, 0}
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if mark[i][j] || grid[i][j] == '#' {
               continue
           }
           mark[i][j] = true
           stack = append(stack, Point{i, j})
           for len(stack) > 0 {
               a := stack[len(stack)-1]
               if vis[a.x][a.y] {
                   // backtracking
                   stack = stack[:len(stack)-1]
                   if len(stack) > 0 {
                       ans = append(ans, Op{'D', a.x, a.y})
                       ans = append(ans, Op{'R', a.x, a.y})
                   }
               } else {
                   vis[a.x][a.y] = true
                   ans = append(ans, Op{'B', a.x, a.y})
                   found := false
                   for p := 0; p < 4; p++ {
                       x := a.x + dx[p]
                       y := a.y + dy[p]
                       if x < 1 || x > n || y < 1 || y > m || mark[x][y] || grid[x][y] == '#' {
                           continue
                       }
                       mark[x][y] = true
                       stack = append(stack, Point{x, y})
                       found = true
                   }
                   if !found {
                       // leaf, change last op to 'R'
                       stack = stack[:len(stack)-1]
                       if len(stack) > 0 {
                           ans[len(ans)-1].Opt = 'R'
                       }
                   }
               }
           }
       }
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for _, op := range ans {
       fmt.Fprintf(writer, "%c %d %d\n", op.Opt, op.x, op.y)
   }
}

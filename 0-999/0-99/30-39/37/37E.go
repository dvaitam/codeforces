package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       row := make([]byte, m)
       // read line
       var line string
       fmt.Fscan(in, &line)
       for j := 0; j < m; j++ {
           row[j] = line[j]
       }
       grid[i] = row
   }
   // directions: up, down, left, right
   dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

   // count connected components of given color
   countComp := func(color byte) int {
       vis := make([][]bool, n)
       for i := range vis {
           vis[i] = make([]bool, m)
       }
       var cnt int
       // simple BFS using slice as queue
       var queue [][2]int
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if !vis[i][j] && grid[i][j] == color {
                   cnt++
                   // start BFS
                   vis[i][j] = true
                   queue = queue[:0]
                   queue = append(queue, [2]int{i, j})
                   for qi := 0; qi < len(queue); qi++ {
                       x, y := queue[qi][0], queue[qi][1]
                       for _, d := range dirs {
                           nx, ny := x+d[0], y+d[1]
                           if nx >= 0 && nx < n && ny >= 0 && ny < m && !vis[nx][ny] && grid[nx][ny] == color {
                               vis[nx][ny] = true
                               queue = append(queue, [2]int{nx, ny})
                           }
                       }
                   }
               }
           }
       }
       return cnt
   }

   bCnt := countComp('B')
   wCnt := countComp('W')
   // minimal strokes: either paint each black component, or paint all black then paint white components
   ans := bCnt
   if wCnt+1 < ans {
       ans = wCnt + 1
   }
   fmt.Println(ans)
}

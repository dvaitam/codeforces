package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n, m, k int
   reader := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }
   freeCount := 0
   var sx, sy int
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' {
               freeCount++
               grid[i][j] = 'X'
               sx, sy = i, j
           }
       }
   }
   toRestore := freeCount - k
   if toRestore > 0 {
       // BFS from last free cell to restore toRestore cells to '.'
       queueX := make([]int, 0, toRestore)
       queueY := make([]int, 0, toRestore)
       head := 0
       grid[sx][sy] = '.'
       toRestore--
       queueX = append(queueX, sx)
       queueY = append(queueY, sy)
       dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
       for head < len(queueX) && toRestore > 0 {
           x := queueX[head]
           y := queueY[head]
           head++
           for _, d := range dirs {
               nx := x + d[0]
               ny := y + d[1]
               if nx < 0 || nx >= n || ny < 0 || ny >= m {
                   continue
               }
               if grid[nx][ny] != 'X' {
                   continue
               }
               grid[nx][ny] = '.'
               toRestore--
               if toRestore == 0 {
                   break
               }
               queueX = append(queueX, nx)
               queueY = append(queueY, ny)
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}

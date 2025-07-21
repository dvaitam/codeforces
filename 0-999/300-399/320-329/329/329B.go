package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var r, c int
   fmt.Fscan(reader, &r, &c)
   grid := make([]string, r)
   for i := 0; i < r; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // find start and exit
   var si, sj, ei, ej int
   for i := 0; i < r; i++ {
       for j := 0; j < c; j++ {
           ch := grid[i][j]
           if ch == 'S' {
               si, sj = i, j
           } else if ch == 'E' {
               ei, ej = i, j
           }
       }
   }
   // BFS from exit to compute dist_exit
   n := r * c
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   idx0 := ei*c + ej
   dist[idx0] = 0
   queue := list.New()
   queue.PushBack(idx0)
   dirs := [4][2]int{{-1,0},{1,0},{0,-1},{0,1}}
   for queue.Len() > 0 {
       e := queue.Remove(queue.Front()).(int)
       x := e / c
       y := e % c
       d := dist[e]
       for _, dir := range dirs {
           ni := x + dir[0]
           nj := y + dir[1]
           if ni < 0 || ni >= r || nj < 0 || nj >= c {
               continue
           }
           if grid[ni][nj] == 'T' {
               continue
           }
           niidx := ni*c + nj
           if dist[niidx] == -1 {
               dist[niidx] = d + 1
               queue.PushBack(niidx)
           }
       }
   }
   startDist := dist[si*c+sj]
   // sum breeders at cells with dist <= startDist
   total := 0
   for i := 0; i < r; i++ {
       for j := 0; j < c; j++ {
           ch := grid[i][j]
           if ch >= '0' && ch <= '9' {
               d := dist[i*c+j]
               if d != -1 && d <= startDist {
                   total += int(ch - '0')
               }
           }
       }
   }
   fmt.Fprintln(writer, total)
}

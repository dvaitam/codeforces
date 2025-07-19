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
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   half := n / 2
   // Prepare grid
   grid := make([][]int, n)
   for i := range grid {
       grid[i] = make([]int, n)
   }
   cnt := 0
   // top-left quadrant: bits end with 3
   for i := 0; i < half; i++ {
       for j := 0; j < half; j++ {
           grid[i][j] = 3 | (cnt << 2)
           cnt++
       }
   }
   cnt = 0
   // top-right quadrant: bits end with 2
   for i := 0; i < half; i++ {
       for j := half; j < n; j++ {
           grid[i][j] = 2 | (cnt << 2)
           cnt++
       }
   }
   cnt = 0
   // bottom-left quadrant: bits end with 0
   for i := half; i < n; i++ {
       for j := 0; j < half; j++ {
           grid[i][j] = cnt << 2
           cnt++
       }
   }
   cnt = 0
   // bottom-right quadrant: bits end with 1
   for i := half; i < n; i++ {
       for j := half; j < n; j++ {
           grid[i][j] = 1 | (cnt << 2)
           cnt++
       }
   }
   // Output grid
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if j > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, grid[i][j])
       }
       writer.WriteByte('\n')
   }
}

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

   var r1, c1, r2, c2, n int
   fmt.Fscan(reader, &r1, &c1, &r2, &c2)
   fmt.Fscan(reader, &n)
   size := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &size[i])
   }

   width := c1
   if c2 > width {
       width = c2
   }
   height := r1 + r2
   grid := make([][]byte, width)
   for i := 0; i < width; i++ {
       row := make([]byte, height)
       for j := 0; j < height; j++ {
           row[j] = '.'
       }
       grid[i] = row
   }

   x, y, dx := 0, 0, 1
   if r1&1 == 1 {
       x = c1 - 1
       dx = -1
   }
   // next position in snake pattern
   next := func() {
       nx := x + dx
       if nx < 0 {
           dx = 1
           x = 0
           y++
       } else if y < r1 && nx >= c1 {
           dx = -1
           x = c1 - 1
           y++
       } else if y >= r1 && nx >= c2 {
           dx = -1
           x = c2 - 1
           y++
       } else {
           x = nx
       }
   }

   for i := 0; i < n; i++ {
       for size[i] > 0 {
           grid[x][y] = byte('a' + i)
           size[i]--
           next()
       }
   }

   fmt.Fprintln(writer, "YES")
   for i := 0; i < width; i++ {
       for j := 0; j < height; j++ {
           writer.WriteByte(grid[i][j])
       }
       writer.WriteByte('\n')
   }
}

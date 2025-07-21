package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       total += a[i]
   }
   // compute y coordinates of points
   y := make([]int, n+1)
   for i := 1; i <= n; i++ {
       if i%2 == 1 {
           y[i] = y[i-1] + a[i-1]
       } else {
           y[i] = y[i-1] - a[i-1]
       }
   }
   // find min and max y among points
   minY, maxY := y[0], y[0]
   for i := 1; i <= n; i++ {
       if y[i] < minY {
           minY = y[i]
       }
       if y[i] > maxY {
           maxY = y[i]
       }
   }
   height := maxY - minY
   width := total
   // prepare grid filled with spaces
   grid := make([][]rune, height)
   for i := 0; i < height; i++ {
       row := make([]rune, width)
       for j := 0; j < width; j++ {
           row[j] = ' '
       }
       grid[i] = row
   }
   // draw segments
   currY := 0
   col := 0
   for i := 1; i <= n; i++ {
       step := a[i-1]
       if i%2 == 1 {
           // up segments
           for k := 0; k < step; k++ {
               y0 := currY + k
               // map y0 to row index: y from maxY-1 down to minY -> rows 0..height-1
               row := maxY - 1 - y0
               grid[row][col] = '/'
               col++
           }
           currY += step
       } else {
           // down segments
           for k := 0; k < step; k++ {
               y0 := currY - 1 - k
               row := maxY - 1 - y0
               grid[row][col] = '\\'
               col++
           }
           currY -= step
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < height; i++ {
       // print each row as string
       for j := 0; j < width; j++ {
           writer.WriteRune(grid[i][j])
       }
       writer.WriteByte('\n')
   }
}

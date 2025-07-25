package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var h, w int
   fmt.Fscan(in, &h, &w)
   grid := make([]string, h)
   for i := 0; i < h; i++ {
       fmt.Fscan(in, &grid[i])
   }
   totalStars := 0
   for i := 0; i < h; i++ {
       for j := 0; j < w; j++ {
           if grid[i][j] == '*' {
               totalStars++
           }
       }
   }
   centerCount := 0
   coveredCells := 0
   for i := 0; i < h; i++ {
       for j := 0; j < w; j++ {
           if grid[i][j] != '*' {
               continue
           }
           up := 0
           for ii := i - 1; ii >= 0 && grid[ii][j] == '*'; ii-- {
               up++
           }
           if up == 0 {
               continue
           }
           down := 0
           for ii := i + 1; ii < h && grid[ii][j] == '*'; ii++ {
               down++
           }
           if down == 0 {
               continue
           }
           left := 0
           for jj := j - 1; jj >= 0 && grid[i][jj] == '*'; jj-- {
               left++
           }
           if left == 0 {
               continue
           }
           right := 0
           for jj := j + 1; jj < w && grid[i][jj] == '*'; jj++ {
               right++
           }
           if right == 0 {
               continue
           }
           centerCount++
           coveredCells = 1 + up + down + left + right
       }
   }
   if centerCount == 1 && coveredCells == totalStars {
       fmt.Fprintln(out, "YES")
   } else {
       fmt.Fprintln(out, "NO")
   }
}

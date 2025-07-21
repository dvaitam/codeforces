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
   // Prepare count of black cells in each 3x3 sub-square (indexed by top-left corner)
   counts := make([][]int, n)
   for i := 0; i < n; i++ {
       counts[i] = make([]int, n)
   }
   ans := -1
   for move := 1; move <= m; move++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       // convert to 0-based
       x--
       y--
       if ans != -1 {
           continue
       }
       // update all 3x3 squares that include (x,y)
       for dx := 0; dx <= 2; dx++ {
           i := x - dx
           if i < 0 || i+2 >= n {
               continue
           }
           for dy := 0; dy <= 2; dy++ {
               j := y - dy
               if j < 0 || j+2 >= n {
                   continue
               }
               counts[i][j]++
               if counts[i][j] == 9 && ans == -1 {
                   ans = move
               }
           }
       }
   }
   if ans == -1 {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, ans)
   }
}

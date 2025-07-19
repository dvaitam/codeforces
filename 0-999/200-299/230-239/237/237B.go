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
   fmt.Fscan(reader, &n)
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, c[i])
       for j := 0; j < c[i]; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }

   type swapOp struct{ x, y, px, py int }
   var swaps []swapOp
   for i := 0; i < n; i++ {
       for j := 0; j < c[i]; j++ {
           mn := a[i][j]
           mi, mj := i, j
           for ii := i; ii < n; ii++ {
               // search from current column index to end of row
               for jj := j; jj < c[ii]; jj++ {
                   if a[ii][jj] < mn {
                       mn = a[ii][jj]
                       mi, mj = ii, jj
                   }
               }
           }
           if mi != i || mj != j {
               swaps = append(swaps, swapOp{i, j, mi, mj})
               a[i][j], a[mi][mj] = a[mi][mj], a[i][j]
           }
       }
   }

   // output
   fmt.Fprintln(writer, len(swaps))
   for _, s := range swaps {
       // convert to 1-based indices
       fmt.Fprintln(writer, s.x+1, s.y+1, s.px+1, s.py+1)
   }
}

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
   grid := make([][]int, n)
   for i := 0; i < n; i++ {
       var row string
       fmt.Fscan(reader, &row)
       grid[i] = make([]int, m)
       for j, c := range row {
           if c == '1' {
               grid[i][j] = 1
           }
       }
   }
   // build column-wise prefix sums over rows
   colPS := make([][]int, n+1)
   for i := range colPS {
       colPS[i] = make([]int, m)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           colPS[i+1][j] = colPS[i][j] + grid[i][j]
       }
   }
   maxP := 0
   // iterate over all pairs of rows
   for r1 := 0; r1 < n; r1++ {
       for r2 := r1; r2 < n; r2++ {
           height := r2 - r1 + 1
           currLen := 0
           for c := 0; c < m; c++ {
               // number of occupied in column c between rows r1..r2
               occ := colPS[r2+1][c] - colPS[r1][c]
               if occ == 0 {
                   currLen++
                   perim := 2 * (currLen + height)
                   if perim > maxP {
                       maxP = perim
                   }
               } else {
                   currLen = 0
               }
           }
       }
   }
   fmt.Fprintln(writer, maxP)
}

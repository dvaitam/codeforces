package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   rowSum := make([]int64, n+1)
   colSum := make([]int64, m+1)
   var c int64
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           fmt.Fscan(reader, &c)
           rowSum[i] += c
           colSum[j] += c
       }
   }
   // Precompute center offsets xj = 4*j - 2, yi = 4*i - 2
   xj := make([]int64, m+1)
   for j := 1; j <= m; j++ {
       xj[j] = int64(4*j - 2)
   }
   yi := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       yi[i] = int64(4*i - 2)
   }
   // Compute optimal li
   bestSy := int64(-1)
   bestLi := 0
   for li := 0; li <= n; li++ {
       var sy int64
       y0 := int64(4 * li)
       for i := 1; i <= n; i++ {
           d := y0 - yi[i]
           sy += rowSum[i] * (d * d)
       }
       if bestSy < 0 || sy < bestSy {
           bestSy = sy
           bestLi = li
       }
   }
   // Compute optimal lj
   bestSx := int64(-1)
   bestLj := 0
   for lj := 0; lj <= m; lj++ {
       var sx int64
       x0 := int64(4 * lj)
       for j := 1; j <= m; j++ {
           d := x0 - xj[j]
           sx += colSum[j] * (d * d)
       }
       if bestSx < 0 || sx < bestSx {
           bestSx = sx
           bestLj = lj
       }
   }
   total := bestSy + bestSx
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total)
   fmt.Fprintf(writer, "%d %d\n", bestLi, bestLj)
}

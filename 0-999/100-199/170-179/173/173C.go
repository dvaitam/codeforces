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
   // Read matrix
   a := make([][]int64, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // Build prefix sum P with size (n+1)x(m+1)
   P := make([][]int64, n+1)
   for i := 0; i <= n; i++ {
       P[i] = make([]int64, m+1)
   }
   for i := 1; i <= n; i++ {
       rowSum := int64(0)
       for j := 1; j <= m; j++ {
           rowSum += a[i-1][j-1]
           P[i][j] = P[i-1][j] + rowSum
       }
   }
   // Find max sum of any odd kxk (k>=3)
   maxSum := int64(-9e18)
   maxK := n
   if m < maxK {
       maxK = m
   }
   // for odd k from 3 to maxK
   for k := 3; k <= maxK; k += 2 {
       kk := k
       for i := 0; i+kk <= n; i++ {
           // iterate rows starting at i, block height kk
           i2 := i + kk
           for j := 0; j+kk <= m; j++ {
               j2 := j + kk
               sum := P[i2][j2] - P[i][j2] - P[i2][j] + P[i][j]
               if sum > maxSum {
                   maxSum = sum
               }
           }
       }
   }
   fmt.Fprint(writer, maxSum)
}

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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   arr := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &arr[i])
   }

   m := n - k + 1
   sums := make([]int64, m+1)
   var window int64
   for i := 1; i <= k; i++ {
       window += arr[i]
   }
   sums[1] = window
   for i := 2; i <= m; i++ {
       window += arr[i+k-1] - arr[i-1]
       sums[i] = window
   }

   bestLeft := make([]int, m+1)
   bestLeft[1] = 1
   for i := 2; i <= m; i++ {
       if sums[i] > sums[bestLeft[i-1]] {
           bestLeft[i] = i
       } else {
           bestLeft[i] = bestLeft[i-1]
       }
   }

   var bestSum int64 = -1
   bestA, bestB := 1, k+1
   for j := k + 1; j <= m; j++ {
       i := bestLeft[j-k]
       total := sums[i] + sums[j]
       if total > bestSum || (total == bestSum && (i < bestA || (i == bestA && j < bestB))) {
           bestSum = total
           bestA = i
           bestB = j
       }
   }

   fmt.Fprintln(writer, bestA, bestB)
}

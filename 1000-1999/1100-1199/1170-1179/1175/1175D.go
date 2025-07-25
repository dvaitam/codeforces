package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Compute suffix sums
   suff := make([]int64, n)
   var sum int64
   for i := n - 1; i >= 0; i-- {
       sum += a[i]
       suff[i] = sum
   }
   // Base cost is all elements multiplied by 1: sum of a = suff[0]
   ans := suff[0]
   if k > 1 {
       // Collect suffix sums starting at positions 1..n-1 (i.e., cut after i-1 yields suff[i])
       vals := make([]int64, n-1)
       for i := 1; i < n; i++ {
           vals[i-1] = suff[i]
       }
       sort.Slice(vals, func(i, j int) bool {
           return vals[i] > vals[j]
       })
       // Add top k-1 values
       need := k - 1
       for i := 0; i < need; i++ {
           ans += vals[i]
       }
   }
   fmt.Fprint(writer, ans)
}

package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int64
   var m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   weights := make([]int64, 0, m)
   for i := 0; i < m; i++ {
       var q, w int64
       fmt.Fscan(reader, &q, &w)
       weights = append(weights, w)
   }
   // maximum number of distinct values k such that k*(k-1)/2 <= n-1
   // solve k^2 - k - 2*(n-1) <= 0
   lim := n - 1
   // compute s = floor(sqrt(1 + 8*lim))
   s := int64(math.Sqrt(float64(1 + 8*lim)))
   k := (1 + s) / 2
   // adjust in case
   for k*(k-1)/2 > lim {
       k--
   }
   // number of coupons to take
   if int64(m) < k {
       k = int64(m)
   }
   sort.Slice(weights, func(i, j int) bool {
       return weights[i] > weights[j]
   })
   var total int64
   for i := int64(0); i < k; i++ {
       total += weights[i]
   }
   fmt.Fprint(writer, total)
}

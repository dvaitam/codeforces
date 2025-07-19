package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   buckets := make([]int, k)
   for i := 0; i < n; i++ {
       var tmp int
       fmt.Fscan(reader, &tmp)
       r := tmp % k
       if r < 0 {
           r += k
       }
       buckets[r]++
   }
   var ans int64
   // handle pairs of remainders
   for i := 1; i <= k/2; i++ {
       if i == k-i {
           // elements that pair among themselves
           ans += int64(2 * (buckets[i] / 2))
       } else {
           ans += int64(2 * min(buckets[i], buckets[k-i]))
       }
   }
   // handle remainder 0
   ans += int64(2 * (buckets[0] / 2))
   fmt.Fprintln(writer, ans)
}

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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // Compute length of longest increasing subsequence (LIS)
   dp := make([]int, 0, n)
   for _, x := range a {
       // find first index where dp[i] >= x
       i := sort.Search(len(dp), func(i int) bool { return dp[i] >= x })
       if i == len(dp) {
           dp = append(dp, x)
       } else {
           dp[i] = x
       }
   }
   fmt.Fprintln(writer, len(dp))
}

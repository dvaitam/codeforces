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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Maximum possible value as per constraints
   const maxA = 100000
   freq := make([]int64, maxA+1)
   var x int
   maxVal := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       if x > maxVal {
           maxVal = x
       }
       freq[x]++
   }

   // dp[i]: max points using values up to i
   dp := make([]int64, maxVal+2)
   dp[0] = 0
   dp[1] = freq[1] * 1
   for i := 2; i <= maxVal; i++ {
       take := dp[i-2] + freq[i]*int64(i)
       if dp[i-1] > take {
           dp[i] = dp[i-1]
       } else {
           dp[i] = take
       }
   }
   fmt.Fprint(writer, dp[maxVal])
}

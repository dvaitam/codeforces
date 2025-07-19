package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   const mV = 1 << 22
   // dp[mask] = max index i where a[i] == mask or subset thereof
   dp := make([]int, mV)
   for i := range dp {
       dp[i] = -1
   }
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       a[i] = x
       if x >= 0 && x < mV {
           dp[x] = i
       }
   }
   // SOS DP over subsets
   for b := 1; b < mV; b <<= 1 {
       for mask := 0; mask < mV; mask++ {
           if mask&b != 0 {
               if dp[mask^b] > dp[mask] {
                   dp[mask] = dp[mask^b]
               }
           }
       }
   }
   fullMask := mV - 1
   for i := 0; i < n; i++ {
       comp := a[i] ^ fullMask
       idx := -1
       if comp >= 0 && comp < mV {
           idx = dp[comp]
       }
       if idx == -1 {
           fmt.Fprint(writer, -1, " ")
       } else {
           fmt.Fprint(writer, a[idx], " ")
       }
   }
   fmt.Fprintln(writer)
}

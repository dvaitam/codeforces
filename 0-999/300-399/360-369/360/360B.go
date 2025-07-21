package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(in, &n, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // binary search on maximum allowed adjacent difference
   low, high := int64(0), int64(2000000005)
   var ans int64
   need := n - k
   for low <= high {
       mid := (low + high) >> 1
       // check if we can keep at least need elements unchanged with max gap mid
       ok := false
       dp := make([]int, n)
       for i := 0; i < n; i++ {
           dp[i] = 1
           for j := 0; j < i; j++ {
               // if keep j and i unchanged, can fill between with max diff mid
               if abs64(a[i]-a[j]) <= mid*int64(i-j) {
                   if dp[j]+1 > dp[i] {
                       dp[i] = dp[j] + 1
                   }
               }
           }
           if dp[i] >= need {
               ok = true
               break
           }
       }
       if ok {
           ans = mid
           high = mid - 1
       } else {
           low = mid + 1
       }
   }
   // print result
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprint(out, ans)
}

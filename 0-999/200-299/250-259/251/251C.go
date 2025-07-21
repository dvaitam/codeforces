package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func lcm(a, b int) int {
   return a / gcd(a, b) * b
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var a, b int64
   var k int
   if _, err := fmt.Fscan(in, &a, &b, &k); err != nil {
       return
   }
   diff := a - b
   if diff <= 0 {
       fmt.Println(0)
       return
   }
   // compute LCM of 2..k
   L := 1
   for x := 2; x <= k; x++ {
       L = lcm(L, x)
   }
   // precompute b mod x
   bmod := make([]int64, k+1)
   for x := 2; x <= k; x++ {
       bmod[x] = b % int64(x)
   }
   // dp[i]: min ops to reduce by exactly i from value b+i
   // need up to L or diff if smaller
   maxI := L
   if diff < int64(L) {
       maxI = int(diff)
   }
   dp := make([]int64, maxI+1)
   dp[0] = 0
   for i := 1; i <= maxI; i++ {
       // subtract 1
       dp[i] = dp[i-1] + 1
       // try mod-x operations
       for x := 2; x <= k; x++ {
           // remainder when at value b+i
           r := (bmod[x] + int64(i%x)) % int64(x)
           if r > 0 && int(r) <= i {
               cand := dp[i-int(r)] + 1
               dp[i] = minInt64(dp[i], cand)
           }
       }
   }
   var ans int64
   if diff <= int64(maxI) {
       ans = dp[int(diff)]
   } else {
       // full chunks
       chunks := diff / int64(L)
       rem := diff % int64(L)
       ans = chunks * dp[L]
       if rem > 0 {
           ans += dp[int(rem)]
       }
   }
   fmt.Println(ans)
}

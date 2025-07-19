package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var l, r int64
   fmt.Fscan(reader, &n, &l, &r)

   // Count numbers in [l, r] with remainder 0,1,2 modulo 3
   a0 := r/3 - (l-1)/3
   a1 := (r+2)/3 - (l+1)/3
   a2 := (r+1)/3 - l/3

   const mod int64 = 1000000007
   // dp0, dp1, dp2 store counts of sequences with sum %3 == 0,1,2
   dp0, dp1, dp2 := int64(1), int64(0), int64(0)
   for i := 0; i < n; i++ {
       ndp0 := (dp0*a0 + dp1*a2 + dp2*a1) % mod
       ndp1 := (dp0*a1 + dp1*a0 + dp2*a2) % mod
       ndp2 := (dp0*a2 + dp1*a1 + dp2*a0) % mod
       dp0, dp1, dp2 = ndp0, ndp1, ndp2
   }

   fmt.Print(dp0)
}

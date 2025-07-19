package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// zfunc computes Z-function of string s
func zfunc(s string) []int {
   n := len(s)
   z := make([]int, n)
   l, r := 0, 0
   for i := 1; i < n; i++ {
       if i <= r {
           z[i] = min(r-i+1, z[i-l])
       }
       for i+z[i] < n && s[z[i]] == s[i+z[i]] {
           z[i]++
       }
       if i+z[i]-1 > r {
           l = i
           r = i + z[i] - 1
       }
   }
   return z
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var aStr, LStr, RStr string
   if _, err := fmt.Fscan(reader, &aStr, &LStr, &RStr); err != nil {
       return
   }
   a := aStr
   n := len(a)
   L := LStr
   R := RStr
   lLen := len(L)
   rLen := len(R)

   // compute Z-functions for L#a and R#a
   t1 := L + "#" + a
   z1 := zfunc(t1)
   t2 := R + "#" + a
   z2 := zfunc(t2)

   dp := make([]int64, n+1)
   sum := make([]int64, n+2)
   dp[n] = 1
   sum[n] = 1

   for i := n - 1; i >= 0; i-- {
       if a[i] == '0' {
           if lLen == 1 && L[0] == '0' {
               dp[i] = dp[i+1]
               sum[i] = (sum[i+1] + dp[i]) % mod
           } else {
               dp[i] = 0
               sum[i] = sum[i+1]
           }
           continue
       }
       if n-i < lLen {
           dp[i] = 0
           sum[i] = sum[i+1]
           continue
       }
       indL := lLen + 1 + i
       indR := rLen + 1 + i
       lf := i + lLen
       // if substring < L, increase lf
       if z1[indL] != lLen && L[z1[indL]] > a[i+z1[indL]] {
           lf++
       }
       var rg int
       if n-i < rLen {
           rg = n
       } else if z2[indR] == rLen || R[z2[indR]] > a[i+z2[indR]] {
           rg = i + rLen
       } else {
           rg = i + rLen - 1
       }
       if rg >= lf {
           add := (sum[lf] - sum[rg+1] + mod) % mod
           dp[i] = add
           sum[i] = (sum[i+1] + add) % mod
       } else {
           dp[i] = 0
           sum[i] = sum[i+1]
       }
   }
   fmt.Fprint(writer, dp[0])
}

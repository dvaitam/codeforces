package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   mod := int64(998244353)
   // Special case k=1
   if k == 1 {
       if n <= 2 {
           fmt.Println(1)
       } else {
           fmt.Println(0)
       }
       return
   }
   // Precompute DP arrays
   f := make([]int64, n+1)
   g := make([]int64, n+1)
   f[1] = int64(k - 1)
   g[1] = int64(k - 2)
   g[0] = 1
   for i := 2; i <= n; i++ {
       g[i] = ((int64(k-2)*g[i-1])%mod + f[i-1]) % mod
       f[i] = int64(k-1) * g[i-1] % mod
   }
   // Read input array
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ans := int64(1)
   process := func(start int) {
       // Build subsequence of positions start, start+2, ...
       b := make([]int, 0, (n+1)/2)
       for i := start; i <= n; i += 2 {
           b = append(b, a[i])
       }
       m := len(b)
       // Check for impossible adjacent known equals
       for i := 0; i+1 < m; i++ {
           if b[i] != -1 && b[i] == b[i+1] {
               fmt.Println(0)
               os.Exit(0)
           }
       }
       temp := 0
       for i, v := range b {
           if v != -1 {
               if temp == i {
                   temp = i + 1
                   continue
               }
               L := i - temp
               if temp == 0 {
                   // Prefix segment of unknowns
                   // use f[L-1] + g[L-1]*(k-1)
                   val := (f[L-1] + g[L-1]*int64(k-1)%mod) % mod
                   ans = ans * val % mod
               } else {
                   // Middle segment between known bounds
                   if b[i] == b[temp-1] {
                       ans = ans * f[L] % mod
                   } else {
                       ans = ans * g[L] % mod
                   }
               }
               temp = i + 1
           }
       }
       if temp != m {
           L := m - temp
           if temp != 0 {
               // Suffix segment
               val := (f[L-1] + g[L-1]*int64(k-1)%mod) % mod
               ans = ans * val % mod
           } else {
               // Entire segment unknown
               if m == 1 {
                   ans = ans * int64(k) % mod
               } else {
                   // k*f[m-2] + k*(k-1)*g[m-2]
                   t := (int64(k)*f[m-2]%mod + int64(k)*int64(k-1)%mod*g[m-2]%mod) % mod
                   ans = ans * t % mod
               }
           }
       }
   }
   // Process odd and even positions
   process(1)
   process(2)
   fmt.Println(ans)
}

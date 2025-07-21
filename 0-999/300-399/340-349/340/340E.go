package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modpow(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
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
   used := make([]bool, n+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] != -1 {
           used[a[i]] = true
       }
   }
   // Count blank positions and forbidden self-matches
   m, k := 0, 0
   for i := 1; i <= n; i++ {
       if a[i-1] == -1 {
           m++
           if !used[i] {
               k++
           }
       }
   }
   // Precompute factorials and inv factorials up to n
   fact := make([]int, n+1)
   invfact := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * i % mod
   }
   invfact[n] = modpow(fact[n], mod-2)
   for i := n; i > 0; i-- {
       invfact[i-1] = invfact[i] * i % mod
   }
   // Inclusion-exclusion over k forbidden matches
   ans := 0
   for i := 0; i <= k; i++ {
       // C(k, i) * (m-i)!
       comb := fact[k] * invfact[i] % mod * invfact[k-i] % mod
       ways := fact[m-i]
       term := comb * ways % mod
       if i%2 == 1 {
           ans = (ans - term + mod) % mod
       } else {
           ans = (ans + term) % mod
       }
   }
   fmt.Fprintln(writer, ans)
}

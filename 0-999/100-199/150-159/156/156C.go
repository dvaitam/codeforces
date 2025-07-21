package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007
const maxN = 3000

var fact [maxN + 1]int64
var invfact [maxN + 1]int64

func modpow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func initFact() {
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invfact[maxN] = modpow(fact[maxN], mod-2)
   for i := maxN; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % mod
   }
}

// comb returns C(n, k)
func comb(n, k int) int64 {
   if k < 0 || k > n || n < 0 {
       return 0
   }
   return fact[n] * invfact[k] % mod * invfact[n-k] % mod
}

func main() {
   initFact()
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       sum := 0
       for j := 0; j < n; j++ {
           sum += int(s[j] - 'a')
       }
       var total int64
       // inclusion-exclusion on upper bound 25 per position
       for k := 0; k <= n; k++ {
           rem := sum - k*26
           if rem < 0 {
               break
           }
           ways := comb(n, k) * comb(rem + n - 1, n - 1) % mod
           if k&1 == 1 {
               total = (total - ways + mod) % mod
           } else {
               total = (total + ways) % mod
           }
       }
       // exclude the word itself
       ans := (total - 1 + mod) % mod
       fmt.Fprintln(writer, ans)
   }
}

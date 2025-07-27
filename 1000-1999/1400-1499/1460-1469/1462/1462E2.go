package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

var fact []int64
var invfact []int64

// fast exponentiation modulo MOD
func modpow(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

// initialize factorials and inverse factorials
func init() {
   const maxN = 200005
   fact = make([]int64, maxN)
   invfact = make([]int64, maxN)
   fact[0] = 1
   for i := 1; i < maxN; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invfact[maxN-1] = modpow(fact[maxN-1], MOD-2)
   for i := maxN - 1; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % MOD
   }
}

// compute C(n, k)
func comb(n, k int) int64 {
   if k < 0 || n < k {
       return 0
   }
   return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m, k int
       fmt.Fscan(reader, &n, &m, &k)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       var ans int64
       l := 0
       for r := 0; r < n; r++ {
           for l < r && a[r]-a[l] > k {
               l++
           }
           // number of elements before r in window is w = r-l
           w := r - l
           // choose m-1 among these to form tuple with a[r] as max
           if w >= m-1 {
               ans = (ans + comb(w, m-1)) % MOD
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD int64 = 1000000009

// modPow computes x^e mod MOD
func modPow(x, e int64) int64 {
   res := int64(1)
   x %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = (res * x) % MOD
       }
       x = (x * x) % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // number of wrong answers
   w := n - m
   // maximum correct without triggering a block of size k
   safeCap := (w + 1) * (k - 1)
   if m <= safeCap {
       // no doubling
       fmt.Println(m % MOD)
       return
   }
   // extra answers beyond safe capacity
   t := m - safeCap
   // number of k-blocks (each triggers a doubling)
   // each block uses k answers; even partial extra cause full block
   b := (t + k - 1) / k
   // Score from b blocks: S_b = 2*k*(2^b - 1)
   pow2b := modPow(2, b)
   blocksScore := (2 * k) % MOD * ((pow2b - 1 + MOD) % MOD) % MOD
   // remaining safe answers give +1 each
   rem := (m - b*k) % MOD
   if rem < 0 {
       rem += MOD
   }
   ans := (blocksScore + rem) % MOD
   fmt.Println(ans)
}

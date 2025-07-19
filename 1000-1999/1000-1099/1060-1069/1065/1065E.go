package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD int64 = 998244353

func modpow(x, y int64) int64 {
   x %= MOD
   var res int64 = 1
   for y > 0 {
       if y&1 == 1 {
           res = res * x % MOD
       }
       x = x * x % MOD
       y >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, A int64
   if _, err := fmt.Fscan(reader, &n, &m, &A); err != nil {
       return
   }
   b := make([]int64, m)
   for i := int64(0); i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   inv2 := modpow(2, MOD-2)
   var res int64 = 1
   var prev int64 = 0
   aMod := A % MOD
   for _, bi := range b {
       diff := bi - prev
       tmp := modpow(aMod, diff)
       res = res * (tmp*(tmp+1)%MOD) % MOD * inv2 % MOD
       prev = bi
   }
   // remaining part
   remain := n - prev*2
   if remain > 0 {
       res = res * modpow(aMod, remain) % MOD
   }
   fmt.Println(res)
}

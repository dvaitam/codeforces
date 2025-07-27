package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modPow computes a^b mod mod
func modPow(a, b int64) int64 {
   res := int64(1)
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
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
   // compute n! mod mod
   var fact int64 = 1
   for i := 1; i <= n; i++ {
       fact = fact * int64(i) % mod
   }
   // compute 2^(n-1) mod mod
   pow2 := modPow(2, int64(n-1))
   // cyclic permutations = n! - 2^(n-1)
   ans := (fact - pow2 + mod) % mod
   fmt.Fprintln(writer, ans)
}

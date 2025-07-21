package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

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

func main() {
   reader := bufio.NewReader(os.Stdin)
   var start, end string
   var k int64
   if _, err := fmt.Fscanln(reader, &start); err != nil {
       return
   }
   if _, err := fmt.Fscanln(reader, &end); err != nil {
       return
   }
   if _, err := fmt.Fscanln(reader, &k); err != nil {
       return
   }
   n := len(start)
   // find shift s such that rotating start by s equals end
   doubled := start + start
   // find all shifts s (0 <= s < n) such that rotating start by s gives end
   cnt := 0
   cnt0 := 0
   // doubled string for substring checks
   for s := 0; s < n; s++ {
       // check if rotated by s equals end
       if doubled[s:s+n] == end {
           cnt++
           if s == 0 {
               cnt0 = 1
           }
       }
   }
   if cnt == 0 {
       fmt.Println(0)
       return
   }
   // count sequences sum mod n == s for each s
   // let t = (n-1)^k mod MOD, p = (-1)^k mod MOD
   t := modpow(int64(n-1), k)
   p := int64(1)
   if k&1 == 1 {
       p = MOD - 1
   }
   invN := modpow(int64(n), MOD-2)
   // Count for s=0: C0 = (t + p*(n-1)) / n
   // Count for s!=0: C1 = (t - p) / n
   C0 := (t + p*int64(n-1)) % MOD * invN % MOD
   C1 := (t - p + MOD) % MOD * invN % MOD
   // total ways
   ways := (int64(cnt0)*C0 + int64(cnt-cnt0)*C1) % MOD
   fmt.Println(ways)
}

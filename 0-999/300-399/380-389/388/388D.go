package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int64
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   // extract bits of k
   // find highest bit
   var L int
   for i := 62; i >= 0; i-- {
       if (k>>i)&1 == 1 {
           L = i
           break
       }
   }
   // precompute powers of two up to L+2
   maxn := L + 2
   pow2 := make([]int64, maxn)
   pow2[0] = 1
   for i := 1; i < maxn; i++ {
       pow2[i] = (pow2[i-1] * 2) % MOD
   }
   // dp[t][tight]
   // tight: 1 if sup prefix equals k prefix so far, 0 if already less
   dp := make([][2]int64, maxn)
   dp[0][1] = 1
   // process bits from L down to 0 with DP on XOR of basis rows
   for pos := L; pos >= 0; pos-- {
       bit := (k >> pos) & 1
       next := make([][2]int64, maxn)
       for t := 0; t < maxn; t++ {
           for tight := 0; tight <= 1; tight++ {
               val := dp[t][tight]
               if val == 0 {
                   continue
               }
               // Case 1: pivot at this bit -> new basis row, XOR bit = 1, multiplier = 1
               if !(tight == 1 && bit == 0) {
                   nt := 0
                   if tight == 1 && bit == 1 {
                       nt = 1
                   }
                   // if tight==0, nt stays 0
                   if t+1 < maxn {
                       next[t+1][nt] = (next[t+1][nt] + val) % MOD
                   }
               }
               // Case 2: free column -> XOR bit 0 or 1
               // parity even: XOR=0, multiplier = 2^{t-1} (or 1 if t==0)
               mult0 := int64(1)
               if t > 0 {
                   mult0 = pow2[t-1]
               }
               // XOR bit = 0
               nt0 := 0
               if tight == 1 {
                   if bit == 0 {
                       nt0 = 1
                   } else {
                       nt0 = 0
                   }
               }
               next[t][nt0] = (next[t][nt0] + val*mult0) % MOD
               // parity odd: XOR=1, multiplier = 2^{t-1}, only if t>0
               if t > 0 {
                   mult1 := pow2[t-1]
                   if !(tight == 1 && bit == 0) {
                       nt1 := 0
                       if tight == 1 && bit == 1 {
                           nt1 = 1
                       }
                       next[t][nt1] = (next[t][nt1] + val*mult1) % MOD
                   }
               }
           }
       }
       dp = next
   }
   // sum all dp[t][tight]
   var ans int64
   for t := 0; t < maxn; t++ {
       ans = (ans + dp[t][0] + dp[t][1]) % MOD
   }
   fmt.Println(ans)
}

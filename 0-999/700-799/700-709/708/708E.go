package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modAdd(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}
func modSub(a, b int) int {
   a -= b
   if a < 0 {
       a += MOD
   }
   return a
}
func modMul(a, b int) int {
   return int(int64(a) * int64(b) % MOD)
}
func modPow(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e&1 != 0 {
           res = modMul(res, x)
       }
       x = modMul(x, x)
       e >>= 1
   }
   return res
}
func modInv(a int) int {
   return modPow(a, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var a, b int
   var k int
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &a, &b)
   fmt.Fscan(in, &k)
   // probability p = a/b mod MOD
   invB := modInv(b)
   p := modMul(a, invB)
   q := modSub(1, p)
   // compute binomial distribution F[i] = C(k,i) p^i q^(k-i) mod MOD for i=0..m-1
   F := make([]int, m)
   // F[0] = q^k
   F[0] = modPow(q, k)
   invQ := modInv(q)
   for i := 1; i < m; i++ {
       // F[i] = F[i-1] * (k-i+1)/i * p/q
       num := (k - (i - 1)) % MOD
       F[i] = modMul(F[i-1], num)
       F[i] = modMul(F[i], modInv(i))
       F[i] = modMul(F[i], p)
       F[i] = modMul(F[i], invQ)
   }
   // prefix sum PF[i] = sum_{j=0..i} F[j]
   PF := make([]int, m)
   PF[0] = F[0]
   for i := 1; i < m; i++ {
       PF[i] = modAdd(PF[i-1], F[i])
   }
   // S0 = sum_{i>=m} F[i] = 1 - PF[m-1]
   S0 := modSub(1, PF[m-1])
   // allocate dp arrays
   size := m + 1
   dp1 := make([]int, size*size)
   dp2 := make([]int, size*size)
   // initialize dp1 for first row: states with l+r <= m-1
   for l := 0; l <= m; l++ {
       var Pl int
       if l < m {
           Pl = F[l]
       } else {
           Pl = S0
       }
       if Pl == 0 {
           continue
       }
       maxR := m - 1 - l
       for r := 0; r <= maxR; r++ {
           var Pr int
           if r < m-l {
               Pr = F[r]
           } else {
               // r == m-l
               // sum_{i>=m-l} F[i]
               if m-l-1 >= 0 {
                   Pr = modSub(1, PF[m-l-1])
               } else {
                   Pr = 1
               }
           }
           dp1[l*size+r] = modMul(Pl, Pr)
       }
   }
   // DP for rows 2..n
   for row := 2; row <= n; row++ {
       // prefix sum S of dp1
       // reuse dp2 for prefix
       for i := 0; i < size; i++ {
           sum := 0
           for j := 0; j < size; j++ {
               sum = modAdd(sum, dp1[i*size+j])
               dp2[i*size+j] = sum
           }
       }
       for j := 0; j < size; j++ {
           for i := 1; i < size; i++ {
               dp2[i*size+j] = modAdd(dp2[i*size+j], dp2[(i-1)*size+j])
           }
       }
       // compute dp2 new values in dp1
       // clear dp1
       for idx := range dp1 {
           dp1[idx] = 0
       }
       for l := 0; l <= m; l++ {
           var Pl int
           if l < m {
               Pl = F[l]
           } else {
               Pl = S0
           }
           if Pl == 0 {
               continue
           }
           for r := 0; r <= m-l-1; r++ {
               var Pr int
               if r < m-l {
                   Pr = F[r]
               } else {
                   if m-l-1 >= 0 {
                       Pr = modSub(1, PF[m-l-1])
                   } else {
                       Pr = 1
                   }
               }
               // sum over dp1 prev where l' <= m-1-r, r' <= m-1-l
               x := m - 1 - r
               y := m - 1 - l
               if x >= m {
                   x = m
               }
               if y >= m {
                   y = m
               }
               sum := dp2[x*size+y]
               if sum != 0 {
                   dp1[l*size+r] = modAdd(dp1[l*size+r], modMul(modMul(Pl, Pr), sum))
               }
           }
       }
   }
   // sum dp1 for final answer
   ans := 0
   for _, v := range dp1 {
       ans = modAdd(ans, v)
   }
   fmt.Println(ans)
}

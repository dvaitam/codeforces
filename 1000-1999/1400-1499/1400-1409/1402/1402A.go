package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int64) int64 {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}
func sub(a, b int64) int64 {
   a -= b
   if a < 0 {
       a += MOD
   }
   return a
}
func mul(a, b int64) int64 {
   return (a * b) % MOD
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   h := make([]int64, n+1)
   w := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &h[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &w[i])
   }
   // prefix sums of widths
   P := make([]int64, n+1)
   SP := make([]int64, n+1)
   SP2 := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       P[i] = add(P[i-1], w[i]%MOD)
       SP[i] = add(SP[i-1], P[i])
       SP2[i] = add(SP2[i-1], mul(P[i], P[i]))
   }
   // monotonic stack for L (prev smaller)
   L := make([]int, n+1)
   stk := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stk) > 0 && h[stk[len(stk)-1]] >= h[i] {
           stk = stk[:len(stk)-1]
       }
       if len(stk) == 0 {
           L[i] = 0
       } else {
           L[i] = stk[len(stk)-1]
       }
       stk = append(stk, i)
   }
   // monotonic stack for R (next smaller or equal)
   R := make([]int, n+1)
   stk = stk[:0]
   for i := n; i >= 1; i-- {
       for len(stk) > 0 && h[stk[len(stk)-1]] > h[i] {
           stk = stk[:len(stk)-1]
       }
       if len(stk) == 0 {
           R[i] = n + 1
       } else {
           R[i] = stk[len(stk)-1]
       }
       stk = append(stk, i)
   }
   inv2 := (MOD + 1) / 2
   var ans int64
   for i := 1; i <= n; i++ {
       // A_i = h_i*(h_i+1)/2 mod
       Ai := mul(h[i]%MOD, add(h[i]%MOD, 1))
       Ai = mul(Ai, inv2)
       li := L[i]
       ri := R[i]
       nl := int64(i - li)
       nr := int64(ri - i)
       // prefix before i
       pref := P[i-1]
       // sum P[j] for j from L..i-1 (j = l-1 for l in L+1..i)
       j0 := li - 1
       if j0 < 0 {
           j0 = 0
       }
       S1 := sub(SP[i-1], SP[j0])
       // sum P[j]^2 j=L..i-1
       S12 := sub(SP2[i-1], SP2[j0])
       // sum x = sum (pref - P[j])
       sum_l_x := sub(mul(nl, pref), S1)
       // sum x^2 = sum (pref-P[j])^2
       // sum x^2 = sum (pref-P[j])^2 = nl*pref^2 -2*pref*S1 + S12
       term1 := mul(nl, mul(pref, pref))
       term2 := mul(sub(0, mul(2, pref)%MOD), S1)
       sum_l_x2 := add(add(term1, term2), S12)
       // sum P[j] j=i..ri-1 (j = r for r in i..R-1)
       S2 := sub(SP[ri-1], SP[i-1])
       // sum P[j]^2 j=i..ri-1
       S22 := sub(SP2[ri-1], SP2[i-1])
       // sum y = sum (P[j] - pref)
       sum_r_y := sub(S2, mul(nr, pref))
       // sum y^2 = sum (P[j]-pref)^2
       // sum y^2 = sum (P[j]-pref)^2 = S22 + nr*pref^2 -2*pref*S2
       term3 := mul(nr, mul(pref, pref))
       term4 := mul(sub(0, mul(2, pref)%MOD), S2)
       sum_r_y2 := add(add(S22, term3), term4)
       // sum t and sum t^2
       sum_t := add(mul(nr, sum_l_x), mul(nl, sum_r_y))
       sum_t2 := add(add(mul(nr, sum_l_x2), mul(nl, sum_r_y2)), mul(2, mul(sum_l_x, sum_r_y)))
       // sum of B = sum (t*(t+1)/2) = (sum_t2 + sum_t)/2
       Ti := mul(inv2, add(sum_t2, sum_t))
       ans = add(ans, mul(Ai, Ti))
   }
   fmt.Println(ans)
}

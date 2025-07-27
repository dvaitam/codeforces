package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const mod = 998244353

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}
func sub(a, b int) int {
   a -= b
   if a < 0 {
       a += mod
   }
   return a
}
func mul(a, b int) int { return int((int64(a) * int64(b)) % mod) }
func powmod(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e&1 != 0 {
           res = mul(res, x)
       }
       x = mul(x, x)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   l := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &l[i], &r[i])
   }
   edges := make([][2]int, m)
   specialMap := make(map[int]int)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--;
       v--;
       edges[i][0], edges[i][1] = u, v
       specialMap[u] = 0
       specialMap[v] = 0
   }
   // map special vertices
   K := 0
   for v := range specialMap {
       specialMap[v] = K
       K++
   }
   // prepare l_r for special
   lsp := make([]int, K)
   rsp := make([]int, K)
   for orig, idx := range specialMap {
       lsp[idx] = l[orig]
       rsp[idx] = r[orig]
   }
   // factorials
   fact := make([]int, n+1)
   invf := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invf[n] = powmod(fact[n], mod-2)
   for i := n; i > 0; i-- {
       invf[i-1] = mul(invf[i], i)
   }
   // c_k
   diff := make([]int, n+2)
   for i := 0; i < n; i++ {
       diff[l[i]]++
       if r[i]+1 <= n {
           diff[r[i]+1]--
       }
   }
   c := make([]int, n+1)
   run := 0
   for k := 0; k <= n; k++ {
       run += diff[k]
       c[k] = run
   }
   // P[s][k]
   P := make([][]int, K+1)
   for s := 0; s <= K; s++ {
       P[s] = make([]int, n+1)
       sum := 0
       for k := 0; k <= n; k++ {
           t := k - s
           var A int
           if t >= 0 {
               rem := c[k] - s
               if rem >= t {
                   A = mul(fact[rem], mul(invf[t], invf[rem-t]))
               }
           }
           sum = add(sum, A)
           P[s][k] = sum
       }
   }
   // edge masks and bounds
   edgeMask := make([]uint64, m)
   Ledge := make([]int, m)
   Redge := make([]int, m)
   for i := 0; i < m; i++ {
       u0, v0 := edges[i][0], edges[i][1]
       ui := specialMap[u0]
       vi := specialMap[v0]
       edgeMask[i] = (1<<ui) | (1<<vi)
       Ledge[i] = lsp[ui]
       if lsp[vi] > Ledge[i] {
           Ledge[i] = lsp[vi]
       }
       Redge[i] = rsp[ui]
       if rsp[vi] < Redge[i] {
           Redge[i] = rsp[vi]
       }
   }
   // DP over masks
   M2 := 1 << m
   vmask := make([]uint64, M2)
   Sarr := make([]int, M2)
   Larr := make([]int, M2)
   Rarr := make([]int, M2)
   Larr[0] = 1
   Rarr[0] = n
   ans := 0
   for mask := 1; mask < M2; mask++ {
       p := bits.TrailingZeros(uint(mask))
       prev := mask & (mask - 1)
       vmask[mask] = vmask[prev] | edgeMask[p]
       Sarr[mask] = bits.OnesCount64(vmask[mask])
       // bounds
       Larr[mask] = Larr[prev]
       if Ledge[p] > Larr[mask] {
           Larr[mask] = Ledge[p]
       }
       Rarr[mask] = Rarr[prev]
       if Redge[p] < Rarr[mask] {
           Rarr[mask] = Redge[p]
       }
   }
   // inclusion-exclusion
   for mask := 0; mask < M2; mask++ {
       s := Sarr[mask]
       Lm := Larr[mask]
       Rm := Rarr[mask]
       if Lm > Rm {
           continue
       }
       sum := P[s][Rm]
       if Lm > 0 {
           sum = sub(sum, P[s][Lm-1])
       }
       if bits.OnesCount(uint(mask))%2 == 0 {
           ans = add(ans, sum)
       } else {
           ans = sub(ans, sum)
       }
   }
   fmt.Fprintln(out, ans)
}

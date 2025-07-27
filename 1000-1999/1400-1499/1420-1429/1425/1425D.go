package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

const mod = 1000000007

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

func mul(a, b int) int {
   return int(int64(a) * int64(b) % mod)
}

func powmod(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e & 1 == 1 {
           res = mul(res, x)
       }
       x = mul(x, x)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M, R int
   fmt.Fscan(in, &N, &M, &R)
   X := make([]int, N)
   Y := make([]int, N)
   B := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &X[i], &Y[i], &B[i])
   }
   // Precompute factorials and invfacts
   fact := make([]int, N+1)
   invfact := make([]int, N+1)
   fact[0] = 1
   for i := 1; i <= N; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invfact[N] = powmod(fact[N], mod-2)
   for i := N; i > 0; i-- {
       invfact[i-1] = mul(invfact[i], i)
   }
   // comb[n] = C(n, M)
   comb := make([]int, N+1)
   for n := M; n <= N; n++ {
       comb[n] = mul(fact[n], mul(invfact[M], invfact[n-M]))
   }
   totalC := comb[N]
   // build bitsets of neighborhoods
   w := (N + 63) / 64
   bs := make([][]uint64, N)
   c := make([]int, N)
   for i := 0; i < N; i++ {
       row := make([]uint64, w)
       xi, yi := X[i], Y[i]
       cnt := 0
       for j := 0; j < N; j++ {
           if max(abs(X[j]-xi), abs(Y[j]-yi)) <= R {
               idx := j / 64; off := uint(j % 64)
               row[idx] |= 1 << off
               cnt++
           }
       }
       bs[i] = row
       c[i] = cnt
   }
   // compute contributions
   ans := 0
   // single i terms
   for i := 0; i < N; i++ {
       cnt := c[i]
       without := 0
       if N-cnt >= M {
           without = comb[N-cnt]
       }
       ways := sub(totalC, without)
       term := mul(mul(B[i], B[i]), ways)
       ans = add(ans, term)
   }
   // pair terms
   for i := 0; i < N; i++ {
       for j := i + 1; j < N; j++ {
           // intersection size
           inter := 0
           for k := 0; k < w; k++ {
               inter += bits.OnesCount64(bs[i][k] & bs[j][k])
           }
           ui := c[i] + c[j] - inter
           // compute inclusion-exclusion
           without_i := 0
           if N-c[i] >= M {
               without_i = comb[N-c[i]]
           }
           without_j := 0
           if N-c[j] >= M {
               without_j = comb[N-c[j]]
           }
           without_ij := 0
           if N-ui >= M {
               without_ij = comb[N-ui]
           }
           ways := totalC
           ways = sub(ways, without_i)
           ways = sub(ways, without_j)
           ways = add(ways, without_ij)
           term := mul(mul(2, mul(B[i], B[j])), ways)
           ans = add(ans, term)
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353
const ROOT = 3

func modAdd(a, b int) int {
   s := a + b
   if s >= MOD {
       s -= MOD
   }
   return s
}
func modSub(a, b int) int {
   s := a - b
   if s < 0 {
       s += MOD
   }
   return s
}
func modMul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
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

// NTT implementation
func ntt(a []int, invert bool) {
   n := len(a)
   // bit-reverse
   j := 0
   for i := 1; i < n; i++ {
       bit := n >> 1
       for ; j&bit != 0; bit >>= 1 {
           j ^= bit
       }
       j |= bit
       if i < j {
           a[i], a[j] = a[j], a[i]
       }
   }
   for length := 2; length <= n; length <<= 1 {
       wlen := modPow(ROOT, (MOD-1)/length)
       if invert {
           wlen = modInv(wlen)
       }
       for i := 0; i < n; i += length {
           w := 1
           half := length >> 1
           for j := 0; j < half; j++ {
               u := a[i+j]
               v := modMul(a[i+j+half], w)
               a[i+j] = modAdd(u, v)
               a[i+j+half] = modSub(u, v)
               w = modMul(w, wlen)
           }
       }
   }
   if invert {
       invN := modInv(n)
       for i := 0; i < n; i++ {
           a[i] = modMul(a[i], invN)
       }
   }
}

func polyMul(a, b []int) []int {
   n := len(a) + len(b) - 1
   sz := 1
   for sz < n {
       sz <<= 1
   }
   fa := make([]int, sz)
   fb := make([]int, sz)
   copy(fa, a)
   copy(fb, b)
   ntt(fa, false)
   ntt(fb, false)
   for i := 0; i < sz; i++ {
       fa[i] = modMul(fa[i], fb[i])
   }
   ntt(fa, true)
   return fa[:n]
}

// polyInv computes inverse series of a, a[0] != 0, up to n terms
func polyInv(a []int, n int) []int {
   res := make([]int, 1)
   res[0] = modInv(a[0])
   for length := 1; length < n; length <<= 1 {
       // compute res up to 2*length
       size := length << 1
       // a0 = a[0:size]
       a0 := make([]int, size)
       copy(a0, a)
       // d = a0 * res
       d := polyMul(a0[:size], res)
       d = d[:size]
       // e = 2 - d
       for i := 0; i < size; i++ {
           d[i] = modSub(0, d[i])
       }
       d[0] = modAdd(d[0], 2)
       // res = res * e
       res = polyMul(res, d)
       if len(res) > size {
           res = res[:size]
       }
   }
   return res[:n]
}

// polySqrt computes sqrt series of a, a[0]=1, up to n terms
func polySqrt(a []int, n int) []int {
   res := make([]int, 1)
   res[0] = 1
   inv2 := modInv(2)
   for length := 1; length < n; length <<= 1 {
       size := length << 1
       // res^2
       tmp2 := polyMul(res, res)
       res2 := make([]int, size)
       copy(res2, tmp2)
       // A0 = a[0:size]
       A0 := make([]int, size)
       copy(A0, a)
       // t0 = A0 + res2
       t0 := make([]int, size)
       for i := 0; i < size; i++ {
           t0[i] = modAdd(A0[i], res2[i])
       }
       // t = t0 * inv(res) up to size
       invRes := polyInv(res, size)
       tmpT := polyMul(t0, invRes)
       t := make([]int, size)
       copy(t, tmpT)
       // res = t * inv2
       for i := 0; i < size; i++ {
           res2[i] = modMul(t[i], inv2)
       }
       res = res2
   }
   return res[:n]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   c := make([]int, n)
   kmin := m + 1
   found := false
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
       if c[i] <= m {
           found = true
           if c[i] < kmin {
               kmin = c[i]
           }
       }
   }
   // if no weight <= m, no trees of positive weight
   if !found {
       for s := 1; s <= m; s++ {
           fmt.Fprintln(writer, 0)
       }
       return
   }
   // build C poly length up to m
   Cpoly := make([]int, m+1)
   for _, v := range c {
       if v <= m {
           Cpoly[v] = 1
       }
   }
   // build A = 1 - 4*Cpoly up to m+kmin
   lim := m + kmin + 1
   A := make([]int, lim)
   A[0] = 1
   for i := 1; i < lim; i++ {
       if i <= m {
           A[i] = modSub(0, modMul(4, Cpoly[i]))
       } else {
           A[i] = 0
       }
   }
   // sqrt A
   S := polySqrt(A, lim)
   // N = 1 - S
   N := make([]int, lim)
   // N[0] = 1 - S[0] = 0
   for i := 1; i < lim; i++ {
       N[i] = modSub(0, S[i])
   }
   // build C' of length lim-kmin
   L := m + 1
   Cp := make([]int, L)
   for i := 0; i < L; i++ {
       idx := i + kmin
       if idx < len(Cpoly) {
           Cp[i] = Cpoly[idx]
       } else {
           Cp[i] = 0
       }
   }
   // invC'
   invCp := polyInv(Cp, L)
   // N'
   Np := make([]int, L)
   for i := 0; i < L; i++ {
       idx := i + kmin
       if idx < len(N) {
           Np[i] = N[idx]
       } else {
           Np[i] = 0
       }
   }
   // F' = Np * invCp
   Fraw := polyMul(Np, invCp)
   inv2 := modInv(2)
   // output f[1..m]
   for s := 1; s <= m; s++ {
       var val int
       if s < len(Fraw) {
           val = modMul(Fraw[s], inv2)
       } else {
           val = 0
       }
       fmt.Fprintln(writer, val)
   }
}

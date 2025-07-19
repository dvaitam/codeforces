package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func modAdd(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}

func modSub(a, b int) int {
   a -= b
   if a < 0 {
       a += mod
   }
   return a
}

func modMul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}

func modPow(a, e int) int {
   res := 1
   x := a % mod
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
   // Fermat's little theorem
   return modPow(a, mod-2)
}

// NTT transforms a in-place; invert indicates inverse transform
func ntt(a []int, invert bool) {
   n := len(a)
   // bit-reverse permutation
   for i, j := 1, 0; i < n; i++ {
       bit := n >> 1
       for j&bit != 0 {
           j ^= bit
           bit >>= 1
       }
       j |= bit
       if i < j {
           a[i], a[j] = a[j], a[i]
       }
   }
   // Cooley-Tukey
   for length := 2; length <= n; length <<= 1 {
       // wlen = primitive root^{(mod-1)/length}
       wlen := modPow(3, (mod-1)/length)
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
       for i := range a {
           a[i] = modMul(a[i], invN)
       }
   }
}

// multiply two polynomials a and b
func multiply(a, b []int) []int {
   need := len(a) + len(b) - 1
   n := 1
   for n < need {
       n <<= 1
   }
   fa := make([]int, n)
   fb := make([]int, n)
   copy(fa, a)
   copy(fb, b)
   ntt(fa, false)
   ntt(fb, false)
   for i := 0; i < n; i++ {
       fa[i] = modMul(fa[i], fb[i])
   }
   ntt(fa, true)
   fa = fa[:need]
   return fa
}

// square polynomial a
func square(a []int) []int {
   return multiply(a, a)
}

var digits [10]bool

// ans returns polynomial for sequences of length n
func ans(n int) []int {
   if n == 0 {
       return []int{1}
   }
   if n%2 == 0 {
       ret := ans(n / 2)
       return square(ret)
   }
   ret := ans(n - 1)
   res := make([]int, len(ret)+10)
   for i, v := range ret {
       if v == 0 {
           continue
       }
       for d := 0; d < 10; d++ {
           if digits[d] {
               res[i+d] = modAdd(res[i+d], v)
           }
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(in, &n, &k)
   for i := 0; i < k; i++ {
       var d int
       fmt.Fscan(in, &d)
       digits[d] = true
   }
   f := ans(n / 2)
   var out int
   for _, v := range f {
       out = modAdd(out, modMul(v, v))
   }
   fmt.Println(out)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353
const root = 3

func modPow(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 != 0 {
           res = int((int64(res) * int64(a)) % mod)
       }
       a = int((int64(a) * int64(a)) % mod)
       e >>= 1
   }
   return res
}

func ntt(a []int, invert bool) {
   n := len(a)
   // bit-reverse permutation
   for i, j := 1, 0; i < n; i++ {
       bit := n >> 1
       for ; j&bit != 0; bit >>= 1 {
           j ^= bit
       }
       j |= bit
       if i < j {
           a[i], a[j] = a[j], a[i]
       }
   }
   // NTT
   for length := 2; length <= n; length <<= 1 {
       // wlen = root^( (mod-1)/length )
       wlen := modPow(root, (mod-1)/length)
       if invert {
           wlen = modPow(wlen, mod-2)
       }
       for i := 0; i < n; i += length {
           w := 1
           half := length >> 1
           for j := 0; j < half; j++ {
               u := a[i+j]
               v := int((int64(a[i+j+half]) * int64(w)) % mod)
               a[i+j] = u + v
               if a[i+j] >= mod {
                   a[i+j] -= mod
               }
               a[i+j+half] = u - v
               if a[i+j+half] < 0 {
                   a[i+j+half] += mod
               }
               w = int((int64(w) * int64(wlen)) % mod)
           }
       }
   }
   if invert {
       invN := modPow(n, mod-2)
       for i := range a {
           a[i] = int((int64(a[i]) * int64(invN)) % mod)
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, x, y int
   fmt.Fscan(in, &n, &x, &y)
   a := make([]int, n+1)
   for i := 0; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Prepare NTT
   // convolution size: >= 2*(x+1)
   size := 1
   for size < 2*(x+1) {
       size <<= 1
   }
   fa := make([]int, size)
   fb := make([]int, size)
   for _, v := range a {
       fa[v] = 1
       fb[x-v] = 1
   }
   ntt(fa, false)
   ntt(fb, false)
   for i := 0; i < size; i++ {
       fa[i] = int((int64(fa[i]) * int64(fb[i])) % mod)
   }
   ntt(fa, true)
   // fa now contains convolution; fa[x + d] > 0 means difference d exists
   // Precompute best lengths
   const maxL = 1000000
   best := make([]int, maxL+1)
   for d := 1; d <= x; d++ {
       idx := x + d
       if idx < len(fa) && fa[idx] > 0 {
           L := 2 * (d + y)
           if L > maxL {
               continue
           }
           for m := L; m <= maxL; m += L {
               if best[m] < L {
                   best[m] = L
               }
           }
       }
   }
   // answer queries
   var q int
   fmt.Fscan(in, &q)
   for i := 0; i < q; i++ {
       var li int
       fmt.Fscan(in, &li)
       if li <= maxL && best[li] > 0 {
           fmt.Fprint(out, best[li])
       } else {
           fmt.Fprint(out, -1)
       }
       if i+1 < q {
           fmt.Fprint(out, " ")
       }
   }
   fmt.Fprintln(out)
}

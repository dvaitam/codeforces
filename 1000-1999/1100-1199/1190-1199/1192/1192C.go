package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}
func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // group words by length
   wordsByLen := make(map[int][]string)
   for i := 0; i < n; i++ {
       var w string
       fmt.Fscan(in, &w)
       l := len(w)
       wordsByLen[l] = append(wordsByLen[l], w)
   }
   total := 0
   for a, ws := range wordsByLen {
       if a < 3 {
           continue
       }
       // build char set of endpoints
       cmap := make(map[byte]int)
       for _, w := range ws {
           c1 := w[0]
           c2 := w[a-1]
           if _, ok := cmap[c1]; !ok {
               cmap[c1] = len(cmap)
           }
           if _, ok := cmap[c2]; !ok {
               cmap[c2] = len(cmap)
           }
       }
       k := len(cmap)
       // build M matrix
       M := make([]int, k*k)
       for _, w := range ws {
           c1 := w[0]
           c2 := w[a-1]
           i := cmap[c1]
           j := cmap[c2]
           // oriented: both directions
           M[i*k+j] = add(M[i*k+j], 1)
           M[j*k+i] = add(M[j*k+i], 1)
       }
       // compute T3 triple tensor: T3[i][j][l] = sum_y M[i][y]*M[j][y]*M[l][y]
       size3 := k * k * k
       T3 := make([]int, size3)
       // tmp for 2d
       tmp2 := make([]int, k*k)
       // v[y] column
       v := make([]int, k)
       for y := 0; y < k; y++ {
           // load column y
           for i := 0; i < k; i++ {
               v[i] = M[i*k+y]
           }
           // tmp2[i][j] = v[i]*v[j]
           for i := 0; i < k; i++ {
               vi := v[i]
               off := i * k
               for j := 0; j < k; j++ {
                   tmp2[off+j] = mul(vi, v[j])
               }
           }
           // accumulate T3[i][j][l] += tmp2[i][j] * v[l]
           for i := 0; i < k; i++ {
               for j := 0; j < k; j++ {
                   t2 := tmp2[i*k+j]
                   if t2 == 0 {
                       continue
                   }
                   base := i*k*k + j*k
                   for l := 0; l < k; l++ {
                       // index = i*k*k + j*k + l
                       idx := base + l
                       T3[idx] = add(T3[idx], mul(t2, v[l]))
                   }
               }
           }
       }
       // sum over x0,x1,x2,x3
       var f int
       // alias T3 for simpler access
       // T3[i][j][l] at idx = i*k*k + j*k + l
       for x0 := 0; x0 < k; x0++ {
           for x1 := 0; x1 < k; x1++ {
               for x2 := 0; x2 < k; x2++ {
                   // precompute t012 = T3[x0][x1][x2]
                   t012 := T3[x0*k*k + x1*k + x2]
                   if t012 == 0 {
                       continue
                   }
                   for x3 := 0; x3 < k; x3++ {
                       // terms: T3[x0][x2][x3] * T3[x0][x1][x3] * T3[x1][x2][x3] * t012
                       t023 := T3[x0*k*k + x2*k + x3]
                       if t023 == 0 {
                           continue
                       }
                       t013 := T3[x0*k*k + x1*k + x3]
                       if t013 == 0 {
                           continue
                       }
                       t123 := T3[x1*k*k + x2*k + x3]
                       if t123 == 0 {
                           continue
                       }
                       // accumulate
                       prod := mul(t012, mul(t023, mul(t013, t123)))
                       f = add(f, prod)
                   }
               }
           }
       }
       total = add(total, f)
   }
   fmt.Println(total)
}

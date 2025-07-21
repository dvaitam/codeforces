package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const mod = 1000000007

// multiply two matrices a and b of size n x n
func matMul(a, b [][]int, n int) [][]int {
   c := make([][]int, n)
   for i := 0; i < n; i++ {
       c[i] = make([]int, n)
       ai := a[i]
       ci := c[i]
       for k := 0; k < n; k++ {
           if ai[k] == 0 {
               continue
           }
           aik := ai[k]
           bk := b[k]
           for j := 0; j < n; j++ {
               ci[j] = (ci[j] + aik*bk[j]) % mod
           }
       }
   }
   return c
}

// multiply vector v (size n) by matrix m (n x n): v * m
func vecMatMul(v []int, m [][]int, n int) []int {
   res := make([]int, n)
   for i := 0; i < n; i++ {
       if v[i] == 0 {
           continue
       }
       vi := v[i]
       mi := m[i]
       for j := 0; j < n; j++ {
           res[j] = (res[j] + vi*mi[j]) % mod
       }
   }
   return res
}

// fast exponentiation of dp * (mat^exp)
func applyPower(dp []int, mat [][]int, exp int) []int {
   n := len(dp)
   // copy base matrix
   base := make([][]int, n)
   for i := 0; i < n; i++ {
       base[i] = make([]int, n)
       copy(base[i], mat[i])
   }
   res := make([]int, n)
   copy(res, dp)
   first := true
   for exp > 0 {
       if exp&1 != 0 {
           if first {
               // res = dp * base
               res = vecMatMul(res, base, n)
               first = false
           } else {
               res = vecMatMul(res, base, n)
           }
       }
       exp >>= 1
       if exp > 0 {
           base = matMul(base, base, n)
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read 7 ints
   w := make([]int, 8)
   for i := 1; i <= 7; i++ {
       s, _ := reader.ReadString(' ')
       if s == "" {
           s, _ = reader.ReadString('\n')
       }
       x, _ := strconv.Atoi(s)
       w[i] = x
   }
   // prepare transition matrices for h=1..7
   T := make([][][]int, 8)
   for h := 1; h <= 7; h++ {
       n := 1 << h
       m := make([][]int, n)
       for a := 0; a < n; a++ {
           m[a] = make([]int, n)
       }
       // for each left mask a, right mask b, count horizontal choices
       // horizontal edges count h-1 positions, mask H from 0..(1<<(h-1))-1
       maxH := 1 << (h - 1)
       for a := 0; a < n; a++ {
           for b := 0; b < n; b++ {
               cnt := 0
               for H := 0; H < maxH; H++ {
                   ok := true
                   for r := 0; r < h; r++ {
                       painted := 0
                       // bottom
                       if r == 0 {
                           painted++
                       } else if H&(1<<(r-1)) != 0 {
                           painted++
                       }
                       // top
                       if r == h-1 {
                           painted++
                       } else if H&(1<<r) != 0 {
                           painted++
                       }
                       // left
                       if a&(1<<r) != 0 {
                           painted++
                       }
                       // right
                       if b&(1<<r) != 0 {
                           painted++
                       }
                       if painted == 4 {
                           ok = false
                           break
                       }
                   }
                   if ok {
                       cnt++
                   }
               }
               m[a][b] = cnt
           }
       }
       T[h] = m
   }
   // dp processing
   var dp []int
   hPrev := 0
   // process segments
   for h := 1; h <= 7; h++ {
       wi := w[h]
       if wi == 0 {
           continue
       }
       // initial dp
       if hPrev == 0 {
           // first segment: dp over masks of size h, only all-ones
           size := 1 << h
           dp = make([]int, size)
           dp[size-1] = 1
       } else {
           // transition between hPrev to h
           size2 := 1 << h
           newDp := make([]int, size2)
           min := hPrev
           if h < hPrev {
               min = h
           }
           onesHi := 0
           if h > hPrev {
               // suffix bits positions [hPrev..h-1]
               onesHi = ((1 << (h - hPrev)) - 1) << hPrev
           }
           onesPrev := 0
           if hPrev > h {
               onesPrev = ((1 << (hPrev - h)) - 1) << h
           }
           for mask1, v := range dp {
               if v == 0 {
                   continue
               }
               // check mask1 valid if hPrev>h: high bits must be 1
               if hPrev > h && (mask1&onesPrev) != onesPrev {
                   continue
               }
               // compute mask2 prefix
               mask2 := mask1 & ((1 << min) - 1)
               mask2 |= onesHi
               newDp[mask2] = (newDp[mask2] + v) % mod
           }
           dp = newDp
       }
       // apply transitions for this segment width wi
       dp = applyPower(dp, T[h], wi)
       hPrev = h
   }
   // after all, right boundary painted: mask must be all ones of hPrev
   ans := dp[(1<<hPrev)-1]
   fmt.Println(ans)
}

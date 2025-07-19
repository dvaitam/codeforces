package main

import (
   "bufio"
   "fmt"
   "math"
   "math/bits"
   "os"
)

var n, m, t, x int
var u, v, w []int
var dis []float64
var p, psum [][]float64
var f, g [][]float64

// fft performs in-place FFT on a, invert indicates inverse transform
func fft(a []complex128, invert bool) {
   n := len(a)
   var logN int
   for (1 << logN) < n {
       logN++
   }
   rev := make([]int, n)
   for i := 0; i < n; i++ {
       rev[i] = rev[i>>1]>>1 | ((i&1) << (logN - 1))
   }
   for i := 0; i < n; i++ {
       if i < rev[i] {
           a[i], a[rev[i]] = a[rev[i]], a[i]
       }
   }
   for lenp := 2; lenp <= n; lenp <<= 1 {
       ang := 2 * math.Pi / float64(lenp)
       if invert {
           ang = -ang
       }
       wlen := complex(math.Cos(ang), math.Sin(ang))
       for i := 0; i < n; i += lenp {
           wcur := complex(1, 0)
           half := lenp >> 1
           for j := 0; j < half; j++ {
               u := a[i+j]
               v := a[i+j+half] * wcur
               a[i+j] = u + v
               a[i+j+half] = u - v
               wcur *= wlen
           }
       }
   }
   if invert {
       invN := complex(1/float64(n), 0)
       for i := 0; i < n; i++ {
           a[i] *= invN
       }
   }
}

// calc computes g[idx][j] for j in [l+ (mid-l)+1 .. r]
func calc(l, r, idx int) {
   mid := (l + r) >> 1
   na := r - l
   nb := mid - l
   nn := na + nb
   // build arrays A and B
   // size N = next power of two >= nn+1
   N := 1 << bits.Len(uint(nn))
   A := make([]complex128, N)
   B := make([]complex128, N)
   for i := 0; i <= na; i++ {
       A[i] = complex(p[idx][i], 0)
   }
   for i := na + 1; i < N; i++ {
       A[i] = 0
   }
   for i := 0; i <= nb; i++ {
       B[i] = complex(f[v[idx]][i+l], 0)
   }
   for i := nb + 1; i < N; i++ {
       B[i] = 0
   }
   fft(A, false)
   fft(B, false)
   for i := 0; i < N; i++ {
       A[i] *= B[i]
   }
   fft(A, true)
   // accumulate results
   for j := mid + 1; j <= r; j++ {
       idxC := j - l
       if idxC >= 0 && idxC < len(A) {
           g[idx][j] += real(A[idxC])
       }
   }
}

// divide and conquer DP on interval [l..r]
func divide(l, r int) {
   if l == r {
       for i := 0; i < m; i++ {
           uu := u[i]
           vv := v[i]
           ww := w[i]
           val := g[i][l] + float64(ww) + psum[i][l+1]*dis[vv]
           if val < f[uu][l] {
               f[uu][l] = val
           }
       }
       return
   }
   mid := (l + r) >> 1
   divide(l, mid)
   for i := 0; i < m; i++ {
       calc(l, r, i)
   }
   divide(mid+1, r)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m, &t, &x)
   u = make([]int, m)
   v = make([]int, m)
   w = make([]int, m)
   p = make([][]float64, m)
   psum = make([][]float64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &u[i], &v[i], &w[i])
       p[i] = make([]float64, t+1)
       psum[i] = make([]float64, t+2)
       for j := 1; j <= t; j++ {
           var pi int
           fmt.Fscan(reader, &pi)
           p[i][j] = float64(pi) / 100000.0
       }
       psum[i][t+1] = 0
       for j := t; j >= 1; j-- {
           psum[i][j] = psum[i][j+1] + p[i][j]
       }
   }
   dis = make([]float64, n+1)
   const INF = 1e18
   for i := 1; i <= n; i++ {
       dis[i] = INF
   }
   dis[n] = 0
   // Bellman-Ford
   for iter := 0; iter < n; iter++ {
       updated := false
       for i := 0; i < m; i++ {
           if dis[v[i]] >= INF {
               continue
           }
           nd := dis[v[i]] + float64(w[i])
           if nd < dis[u[i]] {
               dis[u[i]] = nd
               updated = true
           }
       }
       if !updated {
           break
       }
   }
   // add x
   for i := 1; i <= n; i++ {
       dis[i] += float64(x)
   }
   // init DP arrays
   f = make([][]float64, n+1)
   g = make([][]float64, m)
   for i := 1; i <= n; i++ {
       f[i] = make([]float64, t+1)
       if i != n {
           for j := 0; j <= t; j++ {
               f[i][j] = dis[i]
           }
       }
   }
   for i := 0; i < m; i++ {
       g[i] = make([]float64, t+1)
   }
   divide(0, t)
   // output result
   fmt.Printf("%f\n", f[1][t])
}

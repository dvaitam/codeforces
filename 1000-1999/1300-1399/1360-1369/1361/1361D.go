package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   x := make([]int, n)
   y := make([]int, n)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
       g := gcd(abs(x[i]), abs(y[i]))
       if g != 0 {
           a[i] = x[i] / g
           b[i] = y[i] / g
       }
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   sort.Slice(p, func(i, j int) bool {
       if a[p[i]] != a[p[j]] {
           return a[p[i]] < a[p[j]]
       }
       return b[p[i]] < b[p[j]]
   })

   // collect contributions
   var f []float64
   f = append(f, 0)
   for i := 0; i < n; {
       ai, bi := a[p[i]], b[p[i]]
       if ai == 0 && bi == 0 {
           i++
           continue
       }
       j := i + 1
       for j < n && a[p[j]] == ai && b[p[j]] == bi {
           j++
       }
       m := j - i
       // distances
       dList := make([]float64, m)
       for v := 0; v < m; v++ {
           xi, yi := x[p[i+v]], y[p[i+v]]
           dList[v] = math.Hypot(float64(xi), float64(yi))
       }
       sort.Slice(dList, func(u, v int) bool { return dList[u] > dList[v] })
       for v := 0; v < m; v++ {
           f = append(f, float64(k-2*v-1)*dList[v])
       }
       i = j
   }
   sort.Slice(f, func(i, j int) bool { return f[i] > f[j] })
   var ans float64
   for i := 0; i < k; i++ {
       ans += f[i]
   }

   // dynamic programming on groups
   var f0, f1 float64
   f0 = 0
   f1 = -1e30
   for i := 0; i < n; {
       ai, bi := a[p[i]], b[p[i]]
       if ai == 0 && bi == 0 {
           i++
           continue
       }
       j := i + 1
       for j < n && a[p[j]] == ai && b[p[j]] == bi {
           j++
       }
       m := j - i
       dList := make([]float64, m)
       for v := 0; v < m; v++ {
           xi, yi := x[p[i+v]], y[p[i+v]]
           dList[v] = math.Hypot(float64(xi), float64(yi))
       }
       sort.Slice(dList, func(u, v int) bool { return dList[u] > dList[v] })
       var g0, g1 float64
       for v := 0; v < m; v++ {
           coef := float64(k-2*v-1)
           g0 += coef * dList[v]
           if v < k/2 {
               g1 += coef * dList[v]
           }
           if v >= k/2+n-k {
               g1 += float64(2*(n-v)-k-1) * dList[v]
           }
       }
       if n-k+k/2 > m {
           f1 += g0
           f0 += g0
       } else {
           t := f1 + g0
           u := f0 + g1
           if u > t {
               t = u
           }
           f1 = t
           f0 += g0
       }
       i = j
   }
   if f1 > ans {
       ans = f1
   }
   fmt.Fprintf(writer, "%.10f\n", ans)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func gcd(a, b int) int {
   if b == 0 {
       return a
   }
   return gcd(b, a%b)
}

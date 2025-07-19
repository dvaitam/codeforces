package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MAXN = 2002

var (
   n, h   int
   p, p2  float64
   x       [MAXN]int
   L, R    [MAXN]int
   dp      [MAXN][MAXN][2][2]float64
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// get intersection length of [l, r] and [x0, y0]
func get(l, r, x0, y0 int) float64 {
   if l < x0 {
       l = x0
   }
   if r > y0 {
       r = y0
   }
   return float64(r - l)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &h, &p)
   x[0] = -2000000000
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &x[i])
   }
   x[n+1] = 2000000000
   p2 = 1.0 - p
   sort.Ints(x[1 : n+1])
   L[1] = 1
   for i := 2; i <= n; i++ {
       if x[i]-x[i-1] < h {
           L[i] = L[i-1]
       } else {
           L[i] = i
       }
   }
   R[n] = n
   for i := n - 1; i >= 1; i-- {
       if x[i+1]-x[i] < h {
           R[i] = R[i+1]
       } else {
           R[i] = i
       }
   }
   dp[1][n][0][1] = 1.0
   var res float64
   for i := 1; i <= n; i++ {
       for j := n; j >= i; j-- {
           for a := 0; a < 2; a++ {
               for b := 0; b < 2; b++ {
                   pro := dp[i][j][a][b]
                   if pro == 0 {
                       continue
                   }
                   ll := x[i-1]
                   if a == 1 {
                       ll = x[i-1] + h
                   }
                   rr := x[j+1]
                   if b == 0 {
                       rr = x[j+1] - h
                   }
                   // left pick
                   res += pro * 0.5 * p * get(x[i]-h, x[i], ll, rr)
                   dp[i+1][j][0][b] += pro * 0.5 * p
                   res += pro * 0.5 * p2 * get(x[i], x[R[i]]+h, ll, rr)
                   dp[R[i]+1][j][1][b] += pro * 0.5 * p2
                   // right pick
                   res += pro * 0.5 * p * get(x[L[j]]-h, x[j], ll, rr)
                   dp[i][L[j]-1][a][0] += pro * 0.5 * p
                   res += pro * 0.5 * p2 * get(x[j], x[j]+h, ll, rr)
                   dp[i][j-1][a][1] += pro * 0.5 * p2
               }
           }
       }
   }
   fmt.Fprintf(writer, "%.17f", res)
}

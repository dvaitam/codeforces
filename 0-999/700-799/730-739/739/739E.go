package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   f float64
   c int
}

var n, a, b int
var p1, p2, mi []float64
var dp [][]pair

// por returns the better of two pairs: max f, tiebreak by min c
func por(x, y pair) pair {
   if x.f > y.f || (x.f == y.f && x.c < y.c) {
       return x
   }
   return y
}

// check applies DP with penalty v, returns best (f, count)
func check(v float64) pair {
   // initialize
   negInf := -1e18
   for i := 0; i <= n; i++ {
       for j := 0; j <= a; j++ {
           dp[i][j].f = negInf
           dp[i][j].c = 0
       }
   }
   dp[0][0] = pair{0.0, 0}
   // transitions
   for i := 0; i < n; i++ {
       for j := 0; j <= a; j++ {
           cur := dp[i][j]
           // use poke ball (increment j)
           if j+1 <= a {
               // both balls
               u := pair{cur.f + mi[i] - v, cur.c + 1}
               dp[i+1][j+1] = por(dp[i+1][j+1], u)
               // poke only
               v1 := pair{cur.f + p1[i], cur.c}
               dp[i+1][j+1] = por(dp[i+1][j+1], v1)
           }
           // ultra ball (j stays)
           u2 := pair{cur.f + p2[i] - v, cur.c + 1}
           dp[i+1][j] = por(dp[i+1][j], u2)
           // neither
           v2 := pair{cur.f, cur.c}
           dp[i+1][j] = por(dp[i+1][j], v2)
       }
   }
   return dp[n][a]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &a, &b)
   p1 = make([]float64, n)
   p2 = make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p1[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p2[i])
   }
   mi = make([]float64, n)
   for i := 0; i < n; i++ {
       mi[i] = p1[i] + p2[i] - p1[i]*p2[i]
   }
   dp = make([][]pair, n+1)
   for i := range dp {
       dp[i] = make([]pair, a+1)
   }
   // binary search for penalty
   low, high := 0.0, 1e6
   for it := 0; it < 60; it++ {
       mid := (low + high) / 2
       r := check(mid)
       if r.c > b {
           low = mid
       } else {
           high = mid
       }
   }
   res := check(low).f + float64(b)*low
   fmt.Printf("%.6f\n", res)
}

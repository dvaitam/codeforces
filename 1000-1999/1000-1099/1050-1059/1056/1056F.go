package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const (
   N = 101
   K = 10 * N
   inf = 1e15
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var test int
   if _, err := fmt.Fscan(reader, &test); err != nil {
       return
   }
   // precompute decay factors
   var s [N]float64
   s[0] = 1.0
   for i := 1; i < N; i++ {
       s[i] = s[i-1] * 0.9
   }
   // dp[j][k]: min time for j items and score k
   var dp [N][K]float64
   type pair struct{ a, p int }
   for test > 0 {
       test--
       // reset dp
       for i := 0; i < N; i++ {
           for j := 0; j < K; j++ {
               dp[i][j] = inf
           }
       }
       dp[0][0] = 0.0

       var n int
       var c, t float64
       fmt.Fscan(reader, &n)
       fmt.Fscan(reader, &c, &t)
       tab := make([]pair, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &tab[i].a, &tab[i].p)
       }
       // sort by decreasing a, then decreasing p
       sort.Slice(tab, func(i, j int) bool {
           if tab[i].a != tab[j].a {
               return tab[i].a > tab[j].a
           }
           return tab[i].p > tab[j].p
       })
       // DP over items
       for i := 1; i <= n; i++ {
           a := tab[i-1].a
           p := tab[i-1].p
           for j := i; j > 0; j-- {
               maxK := j * 10
               for k := p; k <= maxK && k < K; k++ {
                   prev := dp[j-1][k-p]
                   if prev+float64(a)/s[j] < dp[j][k] {
                       dp[j][k] = prev + float64(a)/s[j]
                   }
               }
           }
       }
       // compute answer
       ans := 0
       for i := 1; i <= n; i++ {
           for j := 0; j <= i*10 && j < K; j++ {
               x := dp[i][j]
               r := t - float64(i*10)
               if x <= r {
                   if j > ans {
                       ans = j
                   }
               }
               k := math.Sqrt(c * x)
               tr := (k - 1) / c
               if tr > 0 && tr + x/(1+c*tr) <= r {
                   if j > ans {
                       ans = j
                   }
               }
           }
       }
       fmt.Println(ans)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pt struct {
   x        int
   contract int
   y        int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   p := make([]pt, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i].x, &p[i].contract, &p[i].y)
   }
   sort.Slice(p, func(i, j int) bool {
       if p[i].x != p[j].x {
           return p[i].x < p[j].x
       }
       return p[i].y < p[j].y
   })
   dp := make([]int64, n)
   var ans int64
   for i := 0; i < n; i++ {
       dp[i] = 0
       for j := 0; j < i; j++ {
           dx := p[i].x - p[j].x
           sumY := p[i].y + p[j].y
           val := dp[j] + int64(dx)*int64(sumY)*int64(k)
           if val > dp[i] {
               dp[i] = val
           }
       }
       dp[i] -= int64(p[i].contract) * 200
       if dp[i] > ans {
           ans = dp[i]
       }
   }
   // output answer divided by 200 as floating point
   fmt.Printf("%.17f\n", float64(ans)/200.0)
}

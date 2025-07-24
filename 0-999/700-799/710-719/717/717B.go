package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   var c0, c1 int64
   if _, err := fmt.Fscan(in, &n, &c0, &c1); err != nil {
       return
   }
   // dp memoization
   dp := make(map[int64]int64)
   // recursive dp with ternary search on convex split
   var calc func(int64) int64
   calc = func(m int64) int64 {
       if m <= 1 {
           return 0
       }
       if v, ok := dp[m]; ok {
           return v
       }
       // search k in [1, m-1]
       l, r := int64(1), m-1
       var best int64 = -1
       // f evaluates cost for split at k
       var f func(int64) int64
       f = func(k int64) int64 {
           left := calc(k)
           right := calc(m-k)
           return left + right + k*c0 + (m-k)*c1
       }
       // ternary search narrowing
       for r-l > 10 {
           m1 := l + (r-l)/3
           m2 := r - (r-l)/3
           if f(m1) <= f(m2) {
               r = m2 - 1
           } else {
               l = m1 + 1
           }
       }
       // finalize by brute
       for k := l; k <= r; k++ {
           cost := f(k)
           if best < 0 || cost < best {
               best = cost
           }
       }
       dp[m] = best
       return best
   }
   res := calc(n)
   fmt.Println(res)
}

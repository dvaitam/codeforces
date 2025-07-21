package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}


func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   li := make([]int64, n)
   ci := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &li[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ci[i])
   }

   // dp maps gcd value to minimal cost achieving it
   dp := make(map[int64]int64)
   for i := 0; i < n; i++ {
       newdp := make(map[int64]int64)
       // carry over previous states
       for g, cost := range dp {
           if prev, ok := newdp[g]; !ok || cost < prev {
               newdp[g] = cost
           }
       }
       // take current card alone
       if prev, ok := newdp[li[i]]; !ok || ci[i] < prev {
           newdp[li[i]] = ci[i]
       }
       // combine current card with previous states
       for g, cost := range dp {
           ng := gcd(g, li[i])
           nc := cost + ci[i]
           if prev, ok := newdp[ng]; !ok || nc < prev {
               newdp[ng] = nc
           }
       }
       dp = newdp
   }

   if cost, ok := dp[1]; ok {
       fmt.Fprintln(out, cost)
   } else {
       fmt.Fprintln(out, -1)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, k int
       fmt.Fscan(in, &n, &k)
       p := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &p[i])
       }
       // compute prefix sums
       pref := make([]int64, n+1)
       for i := 0; i < n; i++ {
           pref[i+1] = pref[i] + p[i]
       }
       var ans int64
       maxT := n / k
       for tcnt := 1; tcnt <= maxT; tcnt++ {
           // choose x = tcnt * k items: positions [n-x ... n-1]
           // free items are the tcnt cheapest among those: positions L0 .. L0+tcnt-1
           L0 := n - tcnt*k
           R0 := L0 + tcnt - 1
           // sum p[L0..R0]
           sum := pref[R0+1] - pref[L0]
           if sum > ans {
               ans = sum
           }
       }
       fmt.Fprintln(out, ans)
   }
}

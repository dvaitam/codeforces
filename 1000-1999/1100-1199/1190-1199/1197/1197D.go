package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read inputs
   var n, m int
   var k int64
   fmt.Fscan(reader, &n, &m, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   const INF int64 = 9e18
   // mn[c] = minimum of base values for indices l0 with l0 mod m == c
   mn := make([]int64, m)
   for i := 0; i < m; i++ {
       mn[i] = INF
   }
   // initial l0 = 0, pref=0, ql=0, c0=0
   mn[0] = 0

   var ans int64 = 0
   var prefix int64 = 0
   for r := 1; r <= n; r++ {
       prefix += a[r-1]
       cr := r % m
       qr := r / m
       // find minimal base
       var group2 int64 = INF
       for c := cr; c < m; c++ {
           if mn[c] < group2 {
               group2 = mn[c]
           }
       }
       var group1 int64 = INF
       for c := 0; c < cr; c++ {
           if mn[c] < group1 {
               group1 = mn[c]
           }
       }
       best := group2
       if group1 + k < best {
           best = group1 + k
       }
       // compute dp[r]
       cur := prefix - k*int64(qr) - best
       if cur > ans {
           ans = cur
       }
       // update mn for l0 = r
       base := prefix - k*int64(qr)
       if base < mn[cr] {
           mn[cr] = base
       }
   }
   // consider empty subarray -> ans >= 0
   if ans < 0 {
       ans = 0
   }
   fmt.Fprintln(writer, ans)
}

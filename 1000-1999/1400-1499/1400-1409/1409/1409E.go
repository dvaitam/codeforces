package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       var k int64
       fmt.Fscan(in, &n, &k)
       xs := make([]int64, n)
       ys := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &xs[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &ys[i]) // y not used
       }
       sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
       cnt := make([]int, n)
       r := 0
       for i := 0; i < n; i++ {
           for r < n && xs[r] <= xs[i]+k {
               r++
           }
           cnt[i] = r - i
       }
       // suffix max of cnt
       suff := make([]int, n+1)
       for i := n - 1; i >= 0; i-- {
           if cnt[i] > suff[i+1] {
               suff[i] = cnt[i]
           } else {
               suff[i] = suff[i+1]
           }
       }
       ans := 0
       for i := 0; i < n; i++ {
           // second interval starts at i+cnt[i]
           j := i + cnt[i]
           if j > n {
               j = n
           }
           total := cnt[i] + suff[j]
           if total > ans {
               ans = total
           }
       }
       fmt.Fprintln(out, ans)
   }
}

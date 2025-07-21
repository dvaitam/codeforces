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
   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // cnt[k] = number of arrays of size 2^k
   const maxExp = 30
   cnt := make([]int, maxExp+1)
   for i := 0; i < m; i++ {
       var b int
       fmt.Fscan(in, &b)
       if b >= 0 && b <= maxExp {
           cnt[b]++
       }
       // ignore sizes > 2^maxExp (cannot fit)
   }
   // sort clusters by descending size
   sort.Sort(sort.Reverse(sort.IntSlice(a)))
   ans := 0
   // fill each cluster greedily
   for _, cap := range a {
       rem := cap
       // try largest arrays first
       for exp := maxExp; exp >= 0; exp-- {
           size := 1 << exp
           if rem < size || cnt[exp] == 0 {
               continue
           }
           maxk := rem >> exp
           if cnt[exp] <= maxk {
               ans += cnt[exp]
               rem -= cnt[exp] * size
               cnt[exp] = 0
           } else {
               ans += maxk
               cnt[exp] -= maxk
               rem -= maxk * size
           }
       }
   }
   fmt.Fprintln(out, ans)
}

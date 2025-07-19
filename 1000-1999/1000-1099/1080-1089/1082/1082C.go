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
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }

   subs := make([][]int, m)
   for i := 0; i < n; i++ {
       var s int
       var r int
       fmt.Fscan(in, &s, &r)
       subs[s-1] = append(subs[s-1], r)
   }

   ans := make([]int64, n+1)
   for _, cur := range subs {
       if len(cur) == 0 {
           continue
       }
       sort.Slice(cur, func(i, j int) bool { return cur[i] > cur[j] })
       var cum int64
       for j, v := range cur {
           cum += int64(v)
           if cum < 0 {
               break
           }
           ans[j+1] += cum
       }
   }

   var best int64
   for k := 1; k <= n; k++ {
       if ans[k] > best {
           best = ans[k]
       }
   }
   fmt.Fprintln(out, best)
}

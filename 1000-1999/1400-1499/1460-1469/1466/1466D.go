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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(in, &n)
       w := make([]int64, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &w[i])
       }
       deg := make([]int, n+1)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(in, &u, &v)
           deg[u]++
           deg[v]++
       }
       // initial answer for k=1 is total sum of weights
       var total int64
       for i := 1; i <= n; i++ {
           total += w[i]
       }
       // extras: each vertex weight appears deg[i]-1 times
       extras := make([]int64, 0, n)
       for i := 1; i <= n; i++ {
           for j := 1; j < deg[i]; j++ {
               extras = append(extras, w[i])
           }
       }
       // sort extras in descending order
       sort.Slice(extras, func(i, j int) bool {
           return extras[i] > extras[j]
       })
       // build answers
       // we need answers for k = 1 .. n-1
       ans := make([]int64, n)
       ans[1] = total
       // fill for k = 2 .. n-1
       for k := 2; k < n; k++ {
           // extras index k-2
           ans[k] = ans[k-1] + extras[k-2]
       }
       // output
       for k := 1; k < n; k++ {
           if k > 1 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, ans[k])
       }
       out.WriteByte('\n')
   }
}

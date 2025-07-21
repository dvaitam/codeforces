package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   b := make([]int64, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for j := 0; j < m; j++ {
       fmt.Fscan(in, &b[j])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   // prefix sums
   prefA := make([]int64, n+1)
   for i := 0; i < n; i++ {
       prefA[i+1] = prefA[i] + a[i]
   }
   prefB := make([]int64, m+1)
   for i := 0; i < m; i++ {
       prefB[i+1] = prefB[i] + b[i]
   }
   totalB := prefB[m]
   // distinct candidate values
   vs := make([]int64, 0, n+m)
   vs = append(vs, a...)
   vs = append(vs, b...)
   sort.Slice(vs, func(i, j int) bool { return vs[i] < vs[j] })
   uvs := make([]int64, 0, len(vs))
   for _, v := range vs {
       if len(uvs) == 0 || uvs[len(uvs)-1] != v {
           uvs = append(uvs, v)
       }
   }
   // scan through candidates
   var idxA, idxBle int
   var ans int64 = -1
   for _, v := range uvs {
       for idxA < n && a[idxA] < v {
           idxA++
       }
       for idxBle < m && b[idxBle] <= v {
           idxBle++
       }
       cntA := int64(idxA)
       sumA := prefA[idxA]
       cntBgt := int64(m - idxBle)
       sumBle := prefB[idxBle]
       sumBgt := totalB - sumBle
       costA := v*cntA - sumA
       costB := sumBgt - v*cntBgt
       total := costA + costB
       if ans < 0 || total < ans {
           ans = total
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, ans)
}

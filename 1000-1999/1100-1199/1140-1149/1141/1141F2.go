package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type seg struct {
   sum   int
   end   int
   start int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build all subarrays
   v := make([]seg, 0, n*(n+1)/2)
   for i := 0; i < n; i++ {
       sum := 0
       for j := i; j < n; j++ {
           sum += a[j]
           // note: start and end use 1-based indices
           v = append(v, seg{sum: sum, end: j + 1, start: i + 1})
       }
   }
   // sort by sum, then by end, then by start
   sort.Slice(v, func(i, j int) bool {
       if v[i].sum != v[j].sum {
           return v[i].sum < v[j].sum
       }
       if v[i].end != v[j].end {
           return v[i].end < v[j].end
       }
       return v[i].start < v[j].start
   })
   m := len(v)
   dp := make([]int, m)
   b := make([]int, m)
   mx := 0
   // compute dp and b per sum group
   for i := 0; i < m; {
       // find group with same sum
       j := i
       for j+1 < m && v[j+1].sum == v[i].sum {
           j++
       }
       for k := i; k <= j; k++ {
           f := v[k].start
           // earliest end in group
           if v[i].end >= f {
               b[k] = 1
               dp[k] = 1
           } else {
               lo, hi := i, k
               // find largest index lo where v[lo].end < f
               for hi-lo > 1 {
                   md := (lo + hi) >> 1
                   if v[md].end < f {
                       lo = md
                   } else {
                       hi = md
                   }
               }
               if v[hi].end < f {
                   lo = hi
               }
               dp[k] = dp[lo] + 1
               b[k] = dp[lo] + 1
           }
           if k != i && dp[k-1] > dp[k] {
               dp[k] = dp[k-1]
           }
           if dp[k] > mx {
               mx = dp[k]
           }
       }
       i = j + 1
   }
   // reconstruct answer
   original := mx
   ans := make([][2]int, 0, original)
   for i := 0; i < m; {
       j := i
       for j+1 < m && v[j+1].sum == v[i].sum {
           j++
       }
       ok := false
       limit := int(2e9)
       for k := j; k >= i; k-- {
           if b[k] == mx {
               ok = true
               f := v[k].start
               s := v[k].end
               if s < limit {
                   ans = append(ans, [2]int{f, s})
                   mx--
                   limit = f
               }
           }
       }
       i = j + 1
       if ok {
           break
       }
   }
   // output
   fmt.Fprintln(writer, original)
   for i := 0; i < original && i < len(ans); i++ {
       fmt.Fprintf(writer, "%d %d\n", ans[i][0], ans[i][1])
   }
}

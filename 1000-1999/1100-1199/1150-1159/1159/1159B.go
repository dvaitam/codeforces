package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Pair values with indices
   type pair struct{ val, idx int }
   pairs := make([]pair, n)
   for i, v := range a {
       pairs[i] = pair{v, i}
   }
   sort.Slice(pairs, func(i, j int) bool {
       return pairs[i].val > pairs[j].val
   })
   const INF64 = int64(1e18)
   result := INF64
   // S_min: smallest index in set, S_max: largest index
   S_min, S_max := n, -1
   // Process groups of equal values in descending order
   for i := 0; i < n; {
       v := pairs[i].val
       j := i
       for j < n && pairs[j].val == v {
           // insert index
           idx := pairs[j].idx
           if idx < S_min {
               S_min = idx
           }
           if idx > S_max {
               S_max = idx
           }
           j++
       }
       // compute local ratios for this group
       for k := i; k < j; k++ {
           idx := pairs[k].idx
           var local = INF64
           // left distance
           if S_min < idx {
               d := idx - S_min
               r := int64(v) / int64(d)
               if r < local {
                   local = r
               }
           }
           // right distance
           if S_max > idx {
               d := S_max - idx
               r := int64(v) / int64(d)
               if r < local {
                   local = r
               }
           }
           if local < result {
               result = local
           }
       }
       i = j
   }
   if result == INF64 {
       result = 0
   }
   fmt.Fprintln(writer, result)
}

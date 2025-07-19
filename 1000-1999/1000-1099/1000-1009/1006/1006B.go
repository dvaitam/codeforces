package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   v := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }
   // copy and sort descending to get top k
   tmp := make([]int, n)
   copy(tmp, v)
   sort.Slice(tmp, func(i, j int) bool {
       return tmp[i] > tmp[j]
   })
   // sum of top k and multiset of values
   sum := 0
   maxv := make([]int, k)
   for i := 0; i < k; i++ {
       sum += tmp[i]
       maxv[i] = tmp[i]
   }
   // find segments between occurrences
   segments := make([]int, 0, k)
   prev := -1
   var i int
   for i = 0; i < n; i++ {
       if len(maxv) == 0 {
           break
       }
       // look for v[i] in maxv
       for j := 0; j < len(maxv); j++ {
           if v[i] == maxv[j] {
               segments = append(segments, i-prev)
               prev = i
               // remove this occurrence
               maxv = append(maxv[:j], maxv[j+1:]...)
               break
           }
       }
   }
   // add remaining tail to last segment
   if len(segments) > 0 {
       segments[len(segments)-1] += n - i
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, sum)
   for idx, v := range segments {
       if idx > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
}

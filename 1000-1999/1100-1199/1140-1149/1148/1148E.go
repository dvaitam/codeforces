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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       sum += int64(a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
       sum -= int64(b[i])
   }
   if sum != 0 {
       fmt.Fprintln(writer, "NO")
       return
   }
   // pair values with original indices
   type pair struct{ val, id int }
   ap := make([]pair, n)
   for i := 0; i < n; i++ {
       ap[i] = pair{a[i], i + 1}
   }
   sort.Slice(ap, func(i, j int) bool { return ap[i].val < ap[j].val })
   sort.Ints(b)
   // sorted a values
   for i := 0; i < n; i++ {
       a[i] = ap[i].val
   }
   // prepare operations
   type op struct{ i, j, d int }
   ops := make([]op, 0, n)
   j := 1
   for i := 0; i < n; i++ {
       if b[i] < a[i] {
           fmt.Fprintln(writer, "NO")
           return
       }
       for b[i] > a[i] {
           for j < n && b[j] >= a[j] {
               j++
           }
           if j >= n {
               // should not happen if sum == 0
               fmt.Fprintln(writer, "NO")
               return
           }
           // transfer from j to i
           need := b[i] - a[i]
           avail := a[j] - b[j]
           var d int
           if need < avail { d = need } else { d = avail }
           a[i] += d
           a[j] -= d
           ops = append(ops, op{ap[i].id, ap[j].id, d})
       }
   }
   // output result
   fmt.Fprintln(writer, "YES")
   fmt.Fprintln(writer, len(ops))
   for _, o := range ops {
       fmt.Fprintln(writer, o.i, o.j, o.d)
   }
}

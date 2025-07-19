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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   // positions for each value (1..5000)
   const maxV = 5000
   pos := make([][]int, maxV+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       v := a[i]
       if v >= 1 && v <= maxV {
           pos[v] = append(pos[v], i)
       }
   }
   // check feasibility
   if n < k {
       fmt.Fprintln(writer, "NO")
       return
   }
   for v := 1; v <= maxV; v++ {
       if len(pos[v]) > k {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   // build list of all positions in value order
   avail := make([]int, 0, n)
   for v := 1; v <= maxV; v++ {
       for _, idx := range pos[v] {
           avail = append(avail, idx)
       }
   }
   // should be n
   if len(avail) < k {
       fmt.Fprintln(writer, "NO")
       return
   }
   // assign colors cycling through k
   ans := make([]int, n)
   for i, idx := range avail {
       ans[idx] = (i % k) + 1
   }
   fmt.Fprintln(writer, "YES")
   for i, c := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, c)
   }
   fmt.Fprintln(writer)
}

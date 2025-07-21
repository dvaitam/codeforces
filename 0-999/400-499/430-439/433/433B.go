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
   // original prefix sums
   original := make([]int64, n+1)
   v := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &v[i])
       original[i] = original[i-1] + v[i]
   }
   // sorted prefix sums
   u := make([]int64, n)
   for i := 1; i <= n; i++ {
       u[i-1] = v[i]
   }
   sort.Slice(u, func(i, j int) bool { return u[i] < u[j] })
   sorted := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       sorted[i] = sorted[i-1] + u[i-1]
   }

   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var typ, l, r int
       fmt.Fscan(reader, &typ, &l, &r)
       if typ == 1 {
           fmt.Fprintln(writer, original[r]-original[l-1])
       } else {
           fmt.Fprintln(writer, sorted[r]-sorted[l-1])
       }
   }
}

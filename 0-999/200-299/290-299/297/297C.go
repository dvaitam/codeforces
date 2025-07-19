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
   type pair struct { val, idx int }
   a := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].val)
       a[i].idx = i + 1
   }
   sort.Slice(a, func(i, j int) bool {
       return a[i].val < a[j].val
   })
   // compute N as smallest multiple of 3 >= n
   N := ((n + 2) / 3) * 3
   co := make([]int, N)
   for i := 0; i < n; i++ {
       co[i] = a[i].val
   }
   n3 := N / 3
   t := make([][2]int, N)
   // first segment
   for i := 0; i < n3; i++ {
       t[i][0] = co[i]
       t[i][1] = 0
   }
   // second segment
   for i := n3; i < 2*n3; i++ {
       t[i][0] = 0
       t[i][1] = co[i]
   }
   // third segment
   z := n3 - 1
   for i := 2 * n3; i < N; i++ {
       t[i][0] = co[i] - z
       t[i][1] = z
       z--
   }
   // prepare results
   res := make([][2]int, n+1)
   for i := 0; i < n; i++ {
       p := a[i].idx
       res[p][0] = t[i][0]
       res[p][1] = t[i][1]
   }
   // output
   fmt.Fprintln(writer, "YES")
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, res[i][0])
   }
   fmt.Fprintln(writer)
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, res[i][1])
   }
   fmt.Fprintln(writer)
}

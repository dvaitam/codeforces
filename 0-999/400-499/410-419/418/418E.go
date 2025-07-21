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

   var n int
   fmt.Fscan(reader, &n)
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }
   // Precompute row2 and row3 with period 2
   row2 := make([]int, n)
   row3 := make([]int, n)

   recompute := func() {
       cnt := make(map[int]int)
       for i := 0; i < n; i++ {
           v := A[i]
           cnt[v]++
           row2[i] = cnt[v]
       }
       cnt2 := make(map[int]int)
       for i := 0; i < n; i++ {
           v := row2[i]
           cnt2[v]++
           row3[i] = cnt2[v]
       }
   }
   recompute()

   var m int
   fmt.Fscan(reader, &m)
   for qi := 0; qi < m; qi++ {
       var typ, x, y int
       fmt.Fscan(reader, &typ, &x, &y)
       if typ == 1 {
           // update A1
           p := y - 1
           A[p] = x
           recompute()
       } else if typ == 2 {
           row := x
           col := y - 1
           var ans int
           if row <= 1 {
               ans = A[col]
           } else if row%2 == 0 {
               ans = row2[col]
           } else {
               ans = row3[col]
           }
           fmt.Fprintln(writer, ans)
       }
   }
}

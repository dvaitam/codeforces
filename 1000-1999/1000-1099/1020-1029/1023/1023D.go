package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   lo := make([]int, q+1)
   hi := make([]int, q+1)
   for i := 1; i <= q; i++ {
       lo[i] = n
       hi[i] = -1
   }
   for i, v := range a {
       if v > 0 {
           lo[v] = min(lo[v], i)
           hi[v] = max(hi[v], i)
       }
   }
   // ensure q appears
   if lo[q] == n {
       placed := false
       for i := 0; i < n; i++ {
           if a[i] == 0 {
               a[i] = q
               lo[q], hi[q] = i, i
               placed = true
               break
           }
       }
       if !placed {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   // fill segments
   for v := q; v >= 1; v-- {
       if lo[v] > hi[v] {
           continue
       }
       j := lo[v]
       for j < hi[v] {
           if a[j] == 0 {
               a[j] = v
               j++
           } else if a[j] < v {
               fmt.Fprintln(writer, "NO")
               return
           } else if a[j] > v {
               // skip to end of that segment
               j = hi[a[j]]
           } else {
               j++
           }
       }
   }
   // fill remaining zeros
   for i := 1; i < n; i++ {
       if a[i] == 0 {
           a[i] = a[i-1]
       }
   }
   for i := n - 2; i >= 0; i-- {
       if a[i] == 0 {
           a[i] = a[i+1]
       }
   }
   // check any zero
   for i := 0; i < n; i++ {
       if a[i] == 0 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   fmt.Fprintln(writer, "YES")
   for i, v := range a {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}

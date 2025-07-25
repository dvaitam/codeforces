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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   if m == 0 {
       fmt.Fprintln(writer, "YES")
       return
   }
   u0, v0 := a[0], b[0]

   // check if choosing x covers with some y
   check := func(x int) bool {
       // collect edges not covered by x
       idx := make([]int, 0, 16)
       for i := 0; i < m; i++ {
           if a[i] != x && b[i] != x {
               idx = append(idx, i)
           }
       }
       if len(idx) == 0 {
           return true
       }
       // candidate y from first uncovered edge
       i0 := idx[0]
       y1, y2 := a[i0], b[i0]
       // try y1
       ok := true
       for _, i := range idx {
           if a[i] != y1 && b[i] != y1 {
               ok = false
               break
           }
       }
       if ok {
           return true
       }
       // try y2
       ok = true
       for _, i := range idx {
           if a[i] != y2 && b[i] != y2 {
               ok = false
               break
           }
       }
       return ok
   }

   if check(u0) || check(v0) {
       fmt.Fprintln(writer, "YES")
   } else {
       fmt.Fprintln(writer, "NO")
   }
}

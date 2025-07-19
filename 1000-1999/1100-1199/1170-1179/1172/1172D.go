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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Use 1-based indexing
   a := make([]int, n+1)
   b := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   a1 := make([]int, 0, n)
   a2 := make([]int, 0, n)
   a3 := make([]int, 0, n)
   a4 := make([]int, 0, n)

   for i := n; i >= 1; i-- {
       if a[i] == i && b[i] == i {
           continue
       }
       // Replace values equal to i
       ai := a[i]
       bi := b[i]
       for j := 1; j <= i; j++ {
           if a[j] == i {
               a[j] = ai
           }
       }
       for j := 1; j <= i; j++ {
           if b[j] == i {
               b[j] = bi
           }
       }
       a1 = append(a1, ai)
       a2 = append(a2, i)
       a3 = append(a3, i)
       a4 = append(a4, bi)
   }
   m := len(a1)
   fmt.Fprintln(writer, m)
   for k := 0; k < m; k++ {
       fmt.Fprintf(writer, "%d %d %d %d\n", a1[k], a2[k], a3[k], a4[k])
   }
}

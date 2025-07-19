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

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       op := make([][]float64, n+2)
       for i := 0; i < m; i++ {
           var r int
           var p float64
           fmt.Fscan(reader, &r, &p)
           if r <= n {
               op[r] = append(op[r], p)
           }
       }
       sorted := true
       for i := 1; i <= n; i++ {
           if a[i] != i {
               sorted = false
               break
           }
       }
       if sorted {
           fmt.Fprintf(writer, "%.10f\n", 1.0)
           continue
       }
       pos := n
       for pos > 1 && a[pos] == pos {
           pos--
       }
       q := 1.0
       for i := pos; i <= n; i++ {
           for _, p := range op[i] {
               q *= 1 - p
           }
       }
       ans := 1 - q
       fmt.Fprintf(writer, "%.10f\n", ans)
   }
}

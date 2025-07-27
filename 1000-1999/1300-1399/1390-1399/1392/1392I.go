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

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   a := make([]int, n)
   b := make([]int, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   // TODO: implement efficient solution for large n, m
   // Placeholder: output 0 for each query
   var x int
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &x)
       fmt.Fprintln(writer, 0)
   }
}

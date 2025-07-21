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
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   maxA := a[0]
   minA := a[0]
   for _, v := range a {
       if v > maxA {
           maxA = v
       }
       if v < minA {
           minA = v
       }
   }
   minB := b[0]
   for _, v := range b {
       if v < minB {
           minB = v
       }
   }
   v := maxA
   if 2*minA > v {
       v = 2 * minA
   }
   if v < minB {
       fmt.Fprintln(writer, v)
   } else {
       fmt.Fprintln(writer, -1)
   }
}

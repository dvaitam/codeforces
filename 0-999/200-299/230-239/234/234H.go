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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }

   total := n + m
   indexes := make([]int, total)
   flips := make([]int, 0, total)
   ia, ib, im := n-1, m-1, total-1
   c := 0
   for ia >= 0 || ib >= 0 {
       for ia >= 0 && a[ia] == c {
           indexes[im] = ia + 1
           im--
           ia--
       }
       for ib >= 0 && b[ib] == c {
           indexes[im] = ib + 1 + n
           im--
           ib--
       }
       flips = append(flips, im+1)
       if c == 0 {
           c = 1
       } else {
           c = 0
       }
   }

   // output merged order
   for i := 0; i < total; i++ {
       fmt.Fprint(writer, indexes[i])
       if i+1 < total {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)

   // remove last flip marker
   if len(flips) > 0 {
       flips = flips[:len(flips)-1]
   }
   // output flip operations
   fmt.Fprintln(writer, len(flips))
   for i := len(flips) - 1; i >= 0; i-- {
       fmt.Fprint(writer, flips[i])
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}

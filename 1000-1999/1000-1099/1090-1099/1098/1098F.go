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

   var s string
   fmt.Fscan(reader, &s)
   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       t := s[l-1 : r]
       m := len(t)
       if m == 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       // compute Z-function
       z := make([]int, m)
       z[0] = m
       for i := 1; i < m; i++ {
           j := 0
           for i+j < m && t[j] == t[i+j] {
               j++
           }
           z[i] = j
       }
       sum := 0
       for _, v := range z {
           sum += v
       }
       fmt.Fprintln(writer, sum)
   }
}

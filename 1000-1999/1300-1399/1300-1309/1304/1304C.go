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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n int
       var m int64
       fmt.Fscan(reader, &n, &m)
       var last int64 = 0
       low, high := m, m
       ok := true
       for i := 0; i < n; i++ {
           var t, l, h int64
           fmt.Fscan(reader, &t, &l, &h)
           dt := t - last
           low -= dt
           high += dt
           if low < l {
               low = l
           }
           if high > h {
               high = h
           }
           if low > high {
               ok = false
           }
           last = t
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}

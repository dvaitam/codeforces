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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       if k%2 == 1 {
           // odd operations: a[i] = max(a) - a[i]
           mx := a[0]
           for _, v := range a {
               if v > mx {
                   mx = v
               }
           }
           for i, v := range a {
               a[i] = mx - v
           }
       } else {
           // even operations: a[i] = a[i] - min(a)
           mn := a[0]
           for _, v := range a {
               if v < mn {
                   mn = v
               }
           }
           for i, v := range a {
               a[i] = v - mn
           }
       }
       // output result
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprintf(writer, "%d", v)
       }
       writer.WriteByte('\n')
   }
}

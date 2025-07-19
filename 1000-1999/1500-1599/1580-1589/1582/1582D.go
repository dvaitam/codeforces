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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]int64, n)
       if n%2 == 0 {
           for i := 0; i < n; i += 2 {
               b[i] = -a[i+1]
               b[i+1] = a[i]
           }
       } else {
           for i := 0; i < n-3; i += 2 {
               b[i] = -a[i+1]
               b[i+1] = a[i]
           }
           i0 := n - 3
           i1 := n - 2
           i2 := n - 1
           if a[i0]+a[i1] != 0 {
               b[i0] = -a[i2]
               b[i1] = -a[i2]
               b[i2] = a[i0] + a[i1]
           } else if a[i0]+a[i2] != 0 {
               b[i0] = -a[i1]
               b[i1] = a[i0] + a[i2]
               b[i2] = -a[i1]
           } else {
               b[i0] = a[i1] + a[i2]
               b[i1] = -a[i0]
               b[i2] = -a[i0]
           }
       }
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, b[i])
       }
       fmt.Fprint(writer, "\n")
   }
}

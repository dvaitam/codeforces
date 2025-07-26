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

   var n, x, y int
   fmt.Fscan(reader, &n, &x, &y)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   for d := 0; d < n; d++ {
       ok := true
       // check x days before
       start := d - x
       if start < 0 {
           start = 0
       }
       for i := start; i < d; i++ {
           if a[i] < a[d] {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       // check y days after
       end := d + y
       if end >= n {
           end = n - 1
       }
       for j := d + 1; j <= end; j++ {
           if a[j] < a[d] {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, d+1)
           return
       }
   }
}

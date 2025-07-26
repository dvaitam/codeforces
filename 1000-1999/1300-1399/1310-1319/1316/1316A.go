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
   for i := 0; i < t; i++ {
       var n int
       var m int64
       fmt.Fscan(reader, &n, &m)
       var sum int64 = 0
       for j := 0; j < n; j++ {
           var a int64
           fmt.Fscan(reader, &a)
           sum += a
       }
       // Maximum possible score for student 1
       if sum < m {
           fmt.Fprintln(writer, sum)
       } else {
           fmt.Fprintln(writer, m)
       }
   }
}

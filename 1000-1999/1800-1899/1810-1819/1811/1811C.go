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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       b := make([]int, n-1)
       for i := 0; i < n-1; i++ {
           fmt.Fscan(reader, &b[i])
       }
       a := make([]int, n)
       a[0] = b[0]
       for i := 1; i <= n-2; i++ {
           if b[i] < b[i-1] {
               a[i] = b[i]
           } else {
               a[i] = b[i-1]
           }
       }
       a[n-1] = b[n-2]
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprint(a[i]))
       }
       writer.WriteByte('\n')
   }
}

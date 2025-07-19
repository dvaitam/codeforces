package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   as := make([]bool, n)
   var sum int64
   for i := n - 1; i >= 0; i-- {
       if sum > 0 {
           sum -= int64(a[i])
           as[i] = true
       } else {
           sum += int64(a[i])
       }
   }
   if sum > 0 {
       for i := 0; i < n; i++ {
           as[i] = !as[i]
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       if as[i] {
           writer.WriteByte('+')
       } else {
           writer.WriteByte('-')
       }
   }
   writer.WriteByte('\n')
}

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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // output sums: for first element, sum with last; otherwise sum with previous
   for i := 0; i < n; i++ {
       var sum int
       if i == 0 {
           sum = a[0] + a[n-1]
       } else {
           sum = a[i] + a[i-1]
       }
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, sum)
   }
}

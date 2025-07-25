package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var sum1, sum2 int64
   nn := int64(n)
   for i := 0; i < n; i++ {
       x := int64(a[i])
       sum1 += x * (nn - x + 1)
       if i > 0 {
           y := int64(a[i-1])
           mn := y
           mx := x
           if x < y {
               mn = x
               mx = y
           }
           sum2 += mn * (nn - mx + 1)
       }
   }
   res := sum1 - sum2
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, res)
}

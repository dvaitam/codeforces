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
   var l, r int64
   if _, err := fmt.Fscan(reader, &n, &l, &r); err != nil {
       return
   }
   a := make([]int64, n)
   p := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   b := make([]int64, n)
   const inf64 = int64(1<<62 - 1)
   mn := inf64
   mx := -inf64
   for i := 0; i < n; i++ {
       b[i] = a[i] + p[i]
       if b[i] < mn {
           mn = b[i]
       }
       if b[i] > mx {
           mx = b[i]
       }
   }
   if mx-mn <= r-l {
       delta := l - mn
       for i := 0; i < n; i++ {
           b[i] += delta
       }
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteString(" ")
           }
           writer.WriteString(fmt.Sprint(b[i]))
       }
       writer.WriteString("\n")
   } else {
       writer.WriteString("-1\n")
   }
}

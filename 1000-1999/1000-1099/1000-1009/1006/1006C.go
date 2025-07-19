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

   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       a := make([]int64, n+2)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       var suml, sumr, ans int64
       suml = a[1]
       sumr = a[n]
       pl, pr := 1, n
       for pl < pr {
           if suml < sumr {
               pl++
               suml += a[pl]
           } else if sumr < suml {
               pr--
               sumr += a[pr]
           } else {
               ans = suml
               pl++
               suml += a[pl]
               pr--
               sumr += a[pr]
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

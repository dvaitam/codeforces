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
   fmt.Fscan(reader, &n)
   var a, b string
   fmt.Fscan(reader, &a)
   fmt.Fscan(reader, &b)

   // cnt[x][y] counts positions with a[i]=x and b[i]=y
   var cnt [2][2]int64
   for i := 0; i < n; i++ {
       x := int(a[i] - '0')
       y := int(b[i] - '0')
       cnt[x][y]++
   }

   var ans int64
   for i := 0; i < n; i++ {
       xi := int(a[i] - '0')
       yi := int(b[i] - '0')
       other := xi ^ 1
       if yi == 0 {
           ans += cnt[other][1] + cnt[other][0]
       } else {
           ans += cnt[other][0]
       }
   }

   // Each pair counted twice
   fmt.Fprintln(writer, ans/2)
}

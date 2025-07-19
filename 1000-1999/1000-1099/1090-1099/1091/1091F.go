package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       a[i] *= 2
   }
   var s string
   fmt.Fscan(reader, &s)

   var ans, now, nowW, nowG int64
   now = 6
   for i := 0; i < n; i++ {
       switch s[i] {
       case 'W':
           now = 4
           nowW += a[i] / 2
           ans += a[i] * 2
       case 'G':
           tmp := min(nowW, a[i]/2)
           nowW -= tmp
           a[i] -= tmp * 2
           ans += tmp * 4
           nowG += tmp * 2
           nowG += a[i] / 2
           ans += a[i] * 3
       case 'L':
           tmp := min(nowW, a[i]/2)
           nowW -= tmp
           a[i] -= tmp * 2
           ans += tmp * 4
           tmp = min(nowG, a[i]/2)
           nowG -= tmp
           a[i] -= tmp * 2
           ans += tmp * 6
           ans += a[i] * now
       }
   }
   fmt.Fprintln(writer, ans/2)
}

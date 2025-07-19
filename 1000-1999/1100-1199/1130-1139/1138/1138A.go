package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       pre, now, ans := 0, 1, 0
       tmp := a[0]
       for i := 1; i < n; i++ {
           if a[i] == tmp {
               now++
           } else {
               if v := min(pre, now) * 2; v > ans {
                   ans = v
               }
               pre = now
               now = 1
               tmp = a[i]
           }
       }
       if v := min(pre, now) * 2; v > ans {
           ans = v
       }
       fmt.Fprintln(writer, ans)
   }
}

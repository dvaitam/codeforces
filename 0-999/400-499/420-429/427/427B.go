package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, t, c int
   _, _ = fmt.Fscan(reader, &n, &t, &c)
   ans := 0
   cur := 0
   for i := 0; i < n; i++ {
       var x int
       _, _ = fmt.Fscan(reader, &x)
       if x <= t {
           cur++
       } else {
           if cur >= c {
               ans += cur - c + 1
           }
           cur = 0
       }
   }
   if cur >= c {
       ans += cur - c + 1
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}

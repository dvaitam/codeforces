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
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       pos[x] = i
   }
   now := 0
   for i := 0; i < n; i++ {
       var b int
       fmt.Fscan(reader, &b)
       ans := 0
       if pos[b] > now {
           ans = pos[b] - now
           now = pos[b]
       }
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans)
   }
}

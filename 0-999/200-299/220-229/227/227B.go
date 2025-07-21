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
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a >= 1 && a <= n {
           pos[a] = i
       }
   }

   var m int
   fmt.Fscan(reader, &m)
   var vasya, petya int64
   for i := 0; i < m; i++ {
       var b int
       fmt.Fscan(reader, &b)
       if b >= 1 && b <= n {
           p := pos[b]
           vasya += int64(p)
           petya += int64(n - p + 1)
       }
   }
   fmt.Fprintln(writer, vasya, petya)
}

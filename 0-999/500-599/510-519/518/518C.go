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

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   a := make([]int, n+1)
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       pos[a[i]] = i
   }
   var result int64
   for j := 0; j < m; j++ {
       var b int
       fmt.Fscan(reader, &b)
       p := pos[b]
       // screen number = ceil(p/k)
       screen := (p-1)/k + 1
       result += int64(screen)
       if p > 1 {
           // swap with previous
           prevApp := a[p-1]
           a[p-1], a[p] = a[p], a[p-1]
           pos[b] = p - 1
           pos[prevApp] = p
       }
   }
   fmt.Fprint(writer, result)
}

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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   sz := make([]int, n+1)
   S := make([]int, n+1)
   var p int
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       switch {
       case x < m:
           sz[i] = -1
       case x > m:
           sz[i] = 1
       default:
           sz[i] = 0
           p = i
       }
       S[i] = S[i-1] + sz[i]
   }
   offset := n + 1
   cnt := make([]int, 2*n+3)
   for i := 0; i < p; i++ {
       cnt[S[i]+offset]++
   }
   var ans int64
   for i := p; i <= n; i++ {
       idx := S[i] + offset
       ans += int64(cnt[idx] + cnt[idx-1])
   }
   fmt.Fprint(writer, ans)
}

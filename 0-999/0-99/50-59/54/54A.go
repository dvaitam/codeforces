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

   var N, K, C int
   if _, err := fmt.Fscan(reader, &N, &K); err != nil {
       return
   }
   fmt.Fscan(reader, &C)
   holidays := make([]int, C)
   for i := 0; i < C; i++ {
       fmt.Fscan(reader, &holidays[i])
   }

   last := 0
   ans := 0
   for _, h := range holidays {
       // ensure no gap larger than K
       for last+K < h {
           last += K
           ans++
       }
       // gift on holiday
       ans++
       last = h
   }
   // after last holiday, until day N
   for last+K <= N {
       last += K
       ans++
   }
   fmt.Fprint(writer, ans)
}

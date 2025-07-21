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
   // maps for skip-one next pointers
   next1 := make(map[int]int, n)
   next2 := make(map[int]int, n)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       next1[a] = b
       next2[b] = a
   }
   // build even positions: x2, x4, ...
   even := make([]int, 0, n/2)
   cur := 0
   for {
       nxt, ok := next1[cur]
       if !ok || nxt == 0 {
           break
       }
       even = append(even, nxt)
       cur = nxt
   }
   // build odd positions in reverse: x_{n-1}, x_{n-3}, ... x1
   oddRev := make([]int, 0, (n+1)/2)
   cur = 0
   for {
       prv, ok := next2[cur]
       if !ok || prv == 0 {
           break
       }
       oddRev = append(oddRev, prv)
       cur = prv
   }
   // reverse oddRev to get odd positions in forward order: x1, x3, ...
   odd := make([]int, len(oddRev))
   for i, v := range oddRev {
       odd[len(oddRev)-1-i] = v
   }
   // merge odd and even
   res := make([]int, 0, n)
   for i := 0; i < len(odd); i++ {
       res = append(res, odd[i])
       if i < len(even) {
           res = append(res, even[i])
       }
   }
   // output
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(writer, ' ') // use space
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprint(writer, '\n')
}

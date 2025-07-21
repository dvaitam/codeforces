package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   p := make([]int, n)
   e := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i], &e[i])
   }
   const INF = 1 << 60
   ans := INF
   totalMasks := 1 << n
   for mask := 0; mask < totalMasks; mask++ {
       wins := bits.OnesCount(uint(mask))
       effort := 0
       for i := 0; i < n; i++ {
           if (mask>>i)&1 == 1 {
               effort += e[i]
           }
       }
       if effort >= ans {
           continue
       }
       ahead := 0
       for i := 0; i < n; i++ {
           pi := p[i]
           if (mask>>i)&1 == 0 {
               pi++
           }
           if pi > wins || (pi == wins && (mask>>i)&1 == 0) {
               ahead++
           }
       }
       if ahead+1 <= k {
           ans = effort
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}

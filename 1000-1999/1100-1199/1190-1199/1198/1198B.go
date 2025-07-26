package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
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
   // balances and last update times
   val := make([]int, n+1)
   last := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &val[i])
       last[i] = 0
   }
   var q int
   fmt.Fscan(reader, &q)
   // record global payouts
   payouts := make([]int, q+2)
   for i := 1; i <= q; i++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 1 {
           var p, x int
           fmt.Fscan(reader, &p, &x)
           val[p] = x
           last[p] = i
           payouts[i] = 0
       } else {
           var x int
           fmt.Fscan(reader, &x)
           payouts[i] = x
       }
   }
   // build suffix max of payouts
   suf := make([]int, q+3)
   for i := q; i >= 1; i-- {
       suf[i] = max(payouts[i], suf[i+1])
   }
   // compute final balances
   for i := 1; i <= n; i++ {
       // consider payouts after last personal update
       mx := suf[last[i]+1]
       if val[i] < mx {
           val[i] = mx
       }
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, val[i])
   }
   writer.WriteByte('\n')
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, r int
   if _, err := fmt.Fscan(in, &n, &m, &r); err != nil {
       return
   }
   buy := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &buy[i])
   }
   sell := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &sell[i])
   }
   minBuy := buy[0]
   for _, v := range buy {
       if v < minBuy {
           minBuy = v
       }
   }
   maxSell := sell[0]
   for _, v := range sell {
       if v > maxSell {
           maxSell = v
       }
   }
   if maxSell > minBuy {
       shares := r / minBuy
       leftover := r % minBuy
       result := leftover + shares*maxSell
       fmt.Fprintln(out, result)
   } else {
       fmt.Fprintln(out, r)
   }
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var t int64
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   a := make([]int64, n)
   var cost int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       cost += a[i]
   }
   var ans int64
   remd := 0
   for remd != n {
       if cost <= t {
           // number of full rounds we can take
           rounds := t / cost
           ans += rounds * int64(n-remd)
           t -= rounds * cost
       }
       // process individual elements
       for i := 0; i < n; i++ {
           if a[i] == 0 {
               continue
           }
           if a[i] <= t {
               t -= a[i]
               ans++
           } else {
               cost -= a[i]
               a[i] = 0
               remd++
           }
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}

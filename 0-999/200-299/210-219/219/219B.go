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

   var p, d int64
   fmt.Fscan(reader, &p, &d)

   // count existing trailing nines in p
   tmp := p
   bestK := 0
   for tmp%10 == 9 {
       bestK++
       tmp /= 10
   }
   ans := p

   // precompute powers of 10
   pow10 := make([]int64, 20)
   pow10[0] = 1
   for i := 1; i < len(pow10); i++ {
       pow10[i] = pow10[i-1] * 10
   }

   // try to increase number of trailing nines
   for k := bestK + 1; k < len(pow10); k++ {
       mod := pow10[k]
       if p < mod {
           break
       }
       // candidate with last k digits as '9'
       cand := (p/mod)*mod - 1
       if cand < 0 {
           break
       }
       if p-cand <= d {
           bestK = k
           ans = cand
       }
   }

   fmt.Fprintln(writer, ans)
}

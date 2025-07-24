package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var inp string
   if _, err := fmt.Fscan(reader, &inp); err != nil {
       return
   }
   nStr := inp[:len(inp)-1]
   seat := inp[len(inp)-1]
   n, err := strconv.ParseInt(nStr, 10, 64)
   if err != nil {
       return
   }
   // seat order: f(1), e(2), d(3), a(4), b(5), c(6)
   var off int64
   switch seat {
   case 'f': off = 1
   case 'e': off = 2
   case 'd': off = 3
   case 'a': off = 4
   case 'b': off = 5
   case 'c': off = 6
   default:
       off = 0
   }
   // compute block index k (0-based)
   var k int64
   switch n % 4 {
   case 1:
       k = ((n - 1) / 4) * 2
   case 3:
       k = ((n - 3) / 4) * 2
   case 2:
       k = ((n - 2) / 2) + 1
   case 0:
       k = (n - 2) / 2
   }
   // sum of movement times before block k: sum M[i] = 3*k - 2*E, where E = floor((k+1)/2)
   E := (k + 1) / 2
   sumM := 3*k - 2*E
   // total time before serving block k: sum of serve times + movement times
   // serve time per block = 6 seconds
   t0 := 6*k + sumM
   // time when Vasya is served
   result := t0 + off
   fmt.Println(result)
}

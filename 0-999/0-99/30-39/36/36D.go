package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   var k int64
   if _, err := fmt.Fscan(reader, &t, &k); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n, m int64
       fmt.Fscan(reader, &n, &m)
       // compute quotient and remainder for (k+1)
       mod := k + 1
       an, rn := n/mod, n%mod
       bm, rm := m/mod, m%mod
       var win bool
       if rn == rm {
           if rn == k {
               win = true
           } else if rn%2 == 1 {
               win = false
           } else {
               // even remainder: depends on block parity
               win = (an%2 == 1)
           }
       } else {
           // different remainders: parity of difference
           d := rn - rm
           if d < 0 {
               d = -d
           }
           win = (d%2 == 1)
       }
       if win {
           fmt.Print("+")
       } else {
           fmt.Print("-")
       }
       if i < t-1 {
           fmt.Println()
       }
   }
}

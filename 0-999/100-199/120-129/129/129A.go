package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   cntEven, cntOdd, totalParity := 0, 0, 0
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       if x%2 == 0 {
           cntEven++
       } else {
           cntOdd++
           totalParity ^= 1
       }
   }
   if totalParity == 0 {
       fmt.Println(cntEven)
   } else {
       fmt.Println(cntOdd)
   }
}

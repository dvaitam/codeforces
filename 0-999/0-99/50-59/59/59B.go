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
   sum := 0
   minOdd := int(1e9)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       sum += a
       if a%2 != 0 && a < minOdd {
           minOdd = a
       }
   }
   // If total is odd, print sum; else try to subtract smallest odd
   if sum%2 == 1 {
       fmt.Println(sum)
   } else if minOdd < int(1e9) {
       fmt.Println(sum - minOdd)
   } else {
       fmt.Println(0)
   }
}

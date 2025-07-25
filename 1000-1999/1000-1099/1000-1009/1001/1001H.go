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
   // Read n bits (0 or 1)
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       sum ^= (x & 1)
   }
   // Output parity: 1 if odd number of 1s, otherwise 0
   fmt.Println(sum)
}

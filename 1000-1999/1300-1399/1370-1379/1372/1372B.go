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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int64
       fmt.Fscan(reader, &n)
       original := n
       var sp int64 = 0
       // find smallest prime factor
       if n%2 == 0 {
           sp = 2
       } else {
           for i := int64(3); i*i <= n; i += 2 {
               if n%i == 0 {
                   sp = i
                   break
               }
           }
       }
       if sp == 0 {
           // n is prime
           sp = n
       }
       b := original / sp
       a := original - b
       fmt.Fprintln(writer, a, b)
   }
}

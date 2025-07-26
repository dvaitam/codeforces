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
       var n int
       var x int64
       fmt.Fscan(reader, &n, &x)
       var maxa int64
       eq := false
       for i := 0; i < n; i++ {
           var a int64
           fmt.Fscan(reader, &a)
           if a == x {
               eq = true
           }
           if a > maxa {
               maxa = a
           }
       }
       if eq {
           fmt.Fprintln(writer, 1)
       } else {
           // minimum hops: either at least 2, or ceil(x/maxa)
           hops := (x + maxa - 1) / maxa
           if hops < 2 {
               hops = 2
           }
           fmt.Fprintln(writer, hops)
       }
   }
}

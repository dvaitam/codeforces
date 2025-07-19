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
       var n, k int
       fmt.Fscan(reader, &n, &k)
       ans := make([]int, 0, 25)
       // build sequence for k-1
       d := k - 1
       for d > 0 {
           x := (d + 1) / 2
           ans = append(ans, x)
           d -= x
       }
       // fill up using doubling starting from k+1
       d = k + 1
       for len(ans) < 24 {
           ans = append(ans, d)
           d *= 2
       }
       // final element
       ans = append(ans, 2*k+1)

       fmt.Fprintln(writer, len(ans))
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   }
}

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
       // divisible case
       if n%k == 0 {
           fmt.Fprintln(writer, "yEs")
           per := n / k
           for i := 0; i < k; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, per)
           }
           fmt.Fprintln(writer)
           continue
       }
       // try k-1 ones
       k2 := k - 1
       rem1 := n - k2
       if rem1 > 0 && rem1%2 == 1 {
           fmt.Fprintln(writer, "yEs")
           for i := 0; i < k2; i++ {
               fmt.Fprint(writer, 1, " ")
           }
           fmt.Fprintln(writer, rem1)
           continue
       }
       // try k-1 twos
       rem2 := n - 2*k2
       if rem2 > 0 && rem2%2 == 0 {
           fmt.Fprintln(writer, "yEs")
           for i := 0; i < k2; i++ {
               fmt.Fprint(writer, 2, " ")
           }
           fmt.Fprintln(writer, rem2)
           continue
       }
       // no solution
       fmt.Fprintln(writer, "nO")
   }
}

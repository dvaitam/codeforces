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
       fmt.Fscan(reader, &n)
       half := n / 2
       if half%2 == 1 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       sumEven := 0
       sumOdd := 0
       // even numbers descending
       for i := n; i >= 2; i -= 2 {
           fmt.Fprint(writer, i, " ")
           sumEven += i
       }
       // odd numbers descending except last
       for j := n - 3; j >= 1; j -= 2 {
           fmt.Fprint(writer, j, " ")
           sumOdd += j
       }
       // last element to balance
       last := sumEven - sumOdd
       fmt.Fprintln(writer, last)
   }
}

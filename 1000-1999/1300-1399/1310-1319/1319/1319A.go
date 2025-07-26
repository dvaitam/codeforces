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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &r[i])
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   x, y := 0, 0
   for i := 0; i < n; i++ {
       if r[i] == 1 && b[i] == 0 {
           x++
       }
       if r[i] == 0 && b[i] == 1 {
           y++
       }
   }
   if x == 0 {
       fmt.Fprintln(writer, -1)
   } else {
       // minimum k such that k * x > y
       ans := y/x + 1
       fmt.Fprintln(writer, ans)
   }
}

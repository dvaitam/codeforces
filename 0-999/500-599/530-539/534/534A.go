package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   switch n {
   case 1, 2:
       fmt.Fprintln(out, 1)
       fmt.Fprintln(out, 1)
       return
   case 3:
       fmt.Fprintln(out, 2)
       fmt.Fprintln(out, "1 3")
       return
   case 4:
       fmt.Fprintln(out, 4)
       fmt.Fprintln(out, "3 1 4 2")
       return
   default:
       // n >= 5
       fmt.Fprintln(out, n)
       // print odd numbers
       for i := 1; i <= n; i += 2 {
           fmt.Fprint(out, i, " ")
       }
       // print even numbers
       for i := 2; i <= n; i += 2 {
           fmt.Fprint(out, i)
           if i+2 <= n {
               fmt.Fprint(out, " ")
           }
       }
       fmt.Fprintln(out)
   }
}
